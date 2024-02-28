package workload

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	client "sigs.k8s.io/controller-runtime/pkg/client"
)

func getOwnerReference(obj metav1.Object) *metav1.OwnerReference {
	ownerRefs := obj.GetOwnerReferences()
	if ownerRefs == nil || len(ownerRefs) == 0 {
		return nil
	}
	for _, ownerRef := range ownerRefs {
		if ownerRef.Controller == nil || !*ownerRef.Controller {
			return &ownerRef
		}
	}
	return nil
}

func getTopLevelOwnerReference(obj metav1.Object) *metav1.OwnerReference {
	var topLevelOwnerRef *metav1.OwnerReference
	var nextLevelOwnerRef *metav1.OwnerReference

	for {
		nextLevelOwnerRef = getOwnerReference(obj)
		if nextLevelOwnerRef == nil {
			return topLevelOwnerRef
		}
		topLevelOwnerRef = nextLevelOwnerRef
		switch topLevelOwnerRef.Kind {
		case "Deployment":
			deployment, err := client.AppsV1().Deployments(obj.GetNamespace()).Get(topLevelOwnerRef.Name, metav1.GetOptions{})
			if err != nil {
				return topLevelOwnerRef
			}
			obj = deployment
		case "StatefulSet":
			statefulSet, err := client.AppsV1().StatefulSets().Get(topLevelOwnerRef.Name, metav1.GetOptions{})
			if err != nil {
				return topLevelOwnerRef
			}
			obj = statefulSet
		case "Pod":
			pod, err := client.CoreV1().Pods().Get(topLevelOwnerRef.Name, metav1.GetOptions{})
			if err != nil {
				return topLevelOwnerRef
			}
			obj = pod
		case "ReplicaSet":
			replicaSet, err := client.CoreV1().ReplicaSet().Get(topLevelOwnerRef.Name, metav1.GetOptions{})
			if err != nil {
				return topLevelOwnerRef
			}
			obj = replicaSet
		case "CloneSet":
			cloneSet, err := client.CoreV1().CloneSet().Get(topLevelOwnerRef.Name, metav1.GetOptions{})
			if err != nil {
				return topLevelOwnerRef
			}
			obj = cloneSet
		}
	}
}

func ObjectKeyFromObject(obj metav1.Object) {
	panic("unimplemented")
}

// return namespace and name
func GetOwnerReferenceNamespaceName(pod *corev1.Pod) (string, string) {
	topLevelOwnerRef := getOwnerReference(pod.GetObjectMeta())
	// If no top level owner reference, use label
	if topLevelOwnerRef == nil {
		namespace := pod.Namespace
		name, ok := pod.Labels["topOwnerReference"]
		if !ok {
			return "", ""
		}
		return namespace, name
	} else {
		namespace := pod.Namespace
		name := topLevelOwnerRef.Name
		return namespace, name
	}
}
