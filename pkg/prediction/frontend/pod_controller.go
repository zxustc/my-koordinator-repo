/*
 * Used to listen pod create/delete event, and create/delete cr for pods selector by key
 */

package frontend

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	analysisv1alpha1 "github.com/koordinator-sh/koordinator/apis/analysis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	corev1listers "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"

	versioned "github.com/koordinator-sh/koordinator/pkg/client/clientset/versioned"
	analysisinformers "github.com/koordinator-sh/koordinator/pkg/client/informers/externalversions"
	analysislister "github.com/koordinator-sh/koordinator/pkg/client/listers/analysis/v1alpha1"
)

const (
	key = "github.com/set-to-your-own-key"
)

func RegisterPodListener(client client.Client, ctx context.Context) {
	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Fatalf("Failed to create in-cluster config: %v", err)
		return
	}
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Fatalf("Failed to create Kubernetes clientset: %v", err)
		return
	}
	factory := informers.NewSharedInformerFactory(clientset, 10*time.Minute)

	analysisClientSet, err := versioned.NewForConfig(config)
	if err != nil {
		klog.Fatalf("error while init analysis Clientset")
		return
	}

	analysisFactory := analysisinformers.NewSharedInformerFactory(analysisClientSet, 10*time.Minute)

	podInformer := factory.Core().V1().Pods()
	podLister := podInformer.Lister()

	analysisInformer := analysisFactory.Analysis().V1alpha1().Recommendations()
	analysisLister := analysisInformer.Lister()

	podInformer.Informer().AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			ref := pod.OwnerReferences
			if len(ref) == 0 {
				createAnalysisCR(client, analysisLister, pod)
			}
		},
		DeleteFunc: func(obj interface{}) {
			pod := obj.(*corev1.Pod)
			if isLastPod(podLister, pod) {
				deleteAnalysisCR(client, analysisLister, pod)
			}
		},
	})
	if err != nil {
		klog.Errorf("error while register pod listener %+v", err)
	}
	go podInformer.Informer().Run(ctx.Done())
	go analysisInformer.Informer().Run(ctx.Done())
	if !cache.WaitForCacheSync(ctx.Done(),
		podInformer.Informer().HasSynced,
		analysisInformer.Informer().HasSynced) {
		klog.Fatalf("Failed to sync cache for pod controller")
	}
}

func isLastPod(podLister corev1listers.PodLister, pod *corev1.Pod) bool {
	if value, ok := pod.Labels[key]; ok {
		selector := labels.SelectorFromSet(labels.Set{key: value})
		pods, err := podLister.List(selector)
		if err != nil {
			return true
		}
		if len(pods) == 0 {
			return true
		}
	}
	return false
}

func existCR(lister analysislister.RecommendationLister, namespace, name string) bool {
	_, err := lister.Recommendations(namespace).Get(name)
	return err == nil
}

func createAnalysisCR(client client.Client, lister analysislister.RecommendationLister, pod *corev1.Pod) {
	types := analysisv1alpha1.RecommendationPodSelector
	if value, ok := pod.Labels[key]; ok {
		selector := &metav1.LabelSelector{
			MatchLabels: map[string]string{
				key: value,
			},
		}
		workloadName, err := parseWorkloadName(pod)
		if err != nil {
			klog.Errorf("failed to parse pod %s: %+v", pod.Name, err)
			return
		}
		namespace := pod.Namespace
		if existCR(lister, namespace, workloadName) {
			return
		}
		analysisation := &analysisv1alpha1.Recommendation{
			TypeMeta: metav1.TypeMeta{
				Kind:       analysisv1alpha1.RecommendationKind,
				APIVersion: analysisv1alpha1.GroupVersion.String(),
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: namespace,
				Name:      workloadName,
			},
			Spec: analysisv1alpha1.RecommendationSpec{
				Target: analysisv1alpha1.RecommendationTarget{
					Type:        types,
					PodSelector: selector,
					Workload:    nil,
				},
			},
			Status: analysisv1alpha1.RecommendationStatus{},
		}
		err = client.Create(context.TODO(), analysisation)
		if err != nil {
			klog.Infof("failed to create analysisation, pod name %s, err %+v", pod.Name, err)
			return
		}
		klog.Infof("create %s/%s crd success by pod selector.", namespace, workloadName)
	}
}

type TopOwnerReference struct {
	APIVersion         string `json:"apiVersion"`
	Kind               string `json:"kind"`
	Name               string `json:"name"`
	UID                string `json:"uid"`
	Controller         bool   `json:"controller"`
	BlockOwnerDeletion bool   `json:"blockOwnerDeletion"`
}

func parseWorkloadName(pod *corev1.Pod) (string, error) {
	annotationKey := "sd.xhs.com/top-owner-reference"
	if annotation, ok := pod.Annotations[annotationKey]; ok {
		var topOwnerReference TopOwnerReference
		if err := json.Unmarshal([]byte(annotation), &topOwnerReference); err != nil {
			return "", fmt.Errorf("failed to unmarshal annotation: %v", err)
		}
		return topOwnerReference.Name, nil
	}
	return "", fmt.Errorf("annotation %s not found", annotationKey)

}

func deleteAnalysisCR(client client.Client, lister analysislister.RecommendationLister, pod *corev1.Pod) {
	namespace := pod.Namespace
	workloadName, err := parseWorkloadName(pod)
	if err != nil {
		klog.Errorf("no enought info to delete analysisation for %s/%s", namespace, workloadName)
		return
	}
	if !existCR(lister, namespace, workloadName) {
		return
	}
	analysisation := &analysisv1alpha1.Recommendation{
		TypeMeta: metav1.TypeMeta{
			Kind:       analysisv1alpha1.RecommendationKind,
			APIVersion: analysisv1alpha1.GroupVersion.String(),
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      workloadName,
		},
	}
	err = client.Delete(context.TODO(), analysisation)
	if err != nil {
		return
	}
	klog.Infof("delete %s/%s cr sucess by pod selector", namespace, workloadName)
}
