package frontend

import (
	"context"
	"time"

	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/predicate"

	kruisev1alpha1 "github.com/openkruise/kruise-api/apps/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	analysisv1alpha1 "github.com/koordinator-sh/koordinator/apis/analysis/v1alpha1"
	versioned "github.com/koordinator-sh/koordinator/pkg/client/clientset/versioned"
	analysisinformers "github.com/koordinator-sh/koordinator/pkg/client/informers/externalversions"
	analysislister "github.com/koordinator-sh/koordinator/pkg/client/listers/analysis/v1alpha1"
)

type WorkloadReconciler struct {
	Client     client.Client
	Scheme     *runtime.Scheme
	Recorder   record.EventRecorder
	lister     analysislister.RecommendationLister
	namespaces []string
	workloads  []string
}

func AddWorkloadReconciler(mgr ctrl.Manager, namespace []string, enabledWorkload []string) error {
	reconciler := &WorkloadReconciler{
		Client:     mgr.GetClient(),
		Scheme:     mgr.GetScheme(),
		Recorder:   mgr.GetEventRecorderFor("maliang-workload-reconciler"),
		namespaces: make([]string, 0),
		workloads:  make([]string, 0),
	}
	reconciler.namespaces = append(reconciler.namespaces, namespace...)
	klog.Infof("Register recommender for namespaces: %+v", reconciler.namespaces)

	reconciler.workloads = append(reconciler.workloads, enabledWorkload...)
	klog.Infof("Register recommender for workloads: %+v", reconciler.namespaces)

	config, err := rest.InClusterConfig()
	if err != nil {
		klog.Errorf("can not init recommend client for workload controller")
	}
	client, _ := versioned.NewForConfig(config)
	recommendFactory := analysisinformers.NewSharedInformerFactory(client, 10*time.Minute)
	recommendInformer := recommendFactory.Analysis().V1alpha1().Recommendations()
	recommendLister := recommendInformer.Lister()

	go recommendInformer.Informer().Run(context.Background().Done())

	if !cache.WaitForCacheSync(context.Background().Done(),
		recommendInformer.Informer().HasSynced) {
		klog.Fatalf("Failed to sync cache for pod controller")
	}
	reconciler.lister = recommendLister
	return reconciler.SetupWithManager(mgr)
}

func (r *WorkloadReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	var profile bool

	// check if should reconcile?
	profile = false
	for _, n := range r.namespaces {
		if n == "all" {
			profile = true
			break
		}
		if n == req.Namespace {
			profile = true
			break
		}
	}
	if !profile {
		return ctrl.Result{}, nil
	}

	// Try to get the Deployment
	cloneset := &kruisev1alpha1.CloneSet{}
	if err := r.Client.Get(ctx, req.NamespacedName, cloneset); err == nil {
		ownerRef := cloneset.ObjectMeta.OwnerReferences
		if r.IsTopOwnerReference(ownerRef) {
			r.CreateOrUpdateRecommendation(analysisv1alpha1.RecommendationTargetWorkload,
				cloneset.Kind, cloneset.APIVersion, req.NamespacedName, nil, ctx)
		}
		return ctrl.Result{}, nil
	}

	deployment := &appsv1.Deployment{}
	if err := r.Client.Get(ctx, req.NamespacedName, deployment); err == nil {
		ownerRef := deployment.ObjectMeta.OwnerReferences
		if r.IsTopOwnerReference(ownerRef) {
			r.CreateOrUpdateRecommendation(analysisv1alpha1.RecommendationTargetWorkload,
				deployment.Kind, deployment.APIVersion, req.NamespacedName, nil, ctx)
		}
		return ctrl.Result{}, nil
	}

	statefulSet := &appsv1.StatefulSet{}
	if err := r.Client.Get(ctx, req.NamespacedName, statefulSet); err == nil {
		ownerRef := statefulSet.ObjectMeta.OwnerReferences
		if r.IsTopOwnerReference(ownerRef) {
			r.CreateOrUpdateRecommendation(analysisv1alpha1.RecommendationTargetWorkload,
				statefulSet.Kind, statefulSet.APIVersion, req.NamespacedName, nil, ctx)
		}
		return ctrl.Result{}, nil
	}

	daemonset := &appsv1.DaemonSet{}
	if err := r.Client.Get(ctx, req.NamespacedName, daemonset); err == nil {
		ownerRef := daemonset.ObjectMeta.OwnerReferences
		if r.IsTopOwnerReference(ownerRef) {
			r.CreateOrUpdateRecommendation(analysisv1alpha1.RecommendationTargetWorkload,
				daemonset.Kind, daemonset.APIVersion, req.NamespacedName, nil, ctx)
		}
		return ctrl.Result{}, nil
	}

	advancedDaemonset := &kruisev1alpha1.DaemonSet{}
	if err := r.Client.Get(ctx, req.NamespacedName, advancedDaemonset); err == nil {
		ownerRef := advancedDaemonset.ObjectMeta.OwnerReferences
		if r.IsTopOwnerReference(ownerRef) {
			r.CreateOrUpdateRecommendation(analysisv1alpha1.RecommendationTargetWorkload,
				advancedDaemonset.Kind, advancedDaemonset.APIVersion, req.NamespacedName, nil, ctx)
		}
		return ctrl.Result{}, nil
	}

	replicaset := &appsv1.ReplicaSet{}
	if err := r.Client.Get(ctx, req.NamespacedName, replicaset); err == nil {
		ownerRef := replicaset.ObjectMeta.OwnerReferences
		if r.IsTopOwnerReference(ownerRef) {
			r.CreateOrUpdateRecommendation(analysisv1alpha1.RecommendationTargetWorkload,
				replicaset.Kind, replicaset.APIVersion, req.NamespacedName, nil, ctx)
		}
		return ctrl.Result{}, nil
	}

	// If can not get object, it's about delete event.
	r.DeleteRecommendation(ctx, req.NamespacedName)
	return ctrl.Result{}, nil
}

func (r *WorkloadReconciler) CreateOrUpdateRecommendation(
	types analysisv1alpha1.RecommendationTargetType,
	kind, apiversion string, namespacedName types.NamespacedName,
	selector *metav1.LabelSelector, ctx context.Context) error {
	name := namespacedName.Name
	namespace := namespacedName.Namespace
	if r.existCR(namespace, name) {
		return nil
	}
	recommendation := &analysisv1alpha1.Recommendation{
		TypeMeta: metav1.TypeMeta{
			Kind:       analysisv1alpha1.RecommendationKind,
			APIVersion: analysisv1alpha1.GroupVersion.Version,
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
		Spec: analysisv1alpha1.RecommendationSpec{
			Target: analysisv1alpha1.RecommendationTarget{
				Type: types,
				Workload: &analysisv1alpha1.CrossVersionObjectReference{
					Kind:       kind,
					Name:       name,
					Namespace:  namespace,
					APIVersion: apiversion,
				},
				PodSelector: selector,
			},
		},
		Status: analysisv1alpha1.RecommendationStatus{},
	}
	err := r.Client.Create(ctx, recommendation)
	if err == nil {
		klog.Infof("create %s/%s crd success", namespace, name)
		return nil
	}
	// Do not need to update existed resource.
	/*
		if apierrors.IsAlreadyExists(err) {
			existCr := &analysisv1alpha1.Recommendation{}
			if err := r.Client.Get(ctx, client.ObjectKey{
				Namespace: namespace,
				Name:      name,
			}, existCr); err != nil {
				return err
			}
			existCr.Spec = recommendation.Spec
			err := r.Client.Update(ctx, existCr)
			if err != nil {
				klog.Warningf("update crd failed: %+v", err)
				return err
			}
			klog.V(4).Infof("update %s/%s crd success", namespace, name)
		} else {
			klog.Warningf("create crd %s/%sfailed : %+v", namespace, name, err)
		}
	*/
	return nil
}

func (r *WorkloadReconciler) DeleteRecommendation(ctx context.Context, namespacedName types.NamespacedName) error {
	namespace := namespacedName.Namespace
	name := namespacedName.Name
	if !r.existCR(namespace, name) {
		return nil
	}
	cr := &analysisv1alpha1.Recommendation{
		TypeMeta: metav1.TypeMeta{
			Kind:       analysisv1alpha1.RecommendationKind,
			APIVersion: analysisv1alpha1.GroupVersion.Version,
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}
	err := r.Client.Delete(ctx, cr)
	if err != nil {
		klog.Warningf("delete %s/%s cr failed: %+v", namespace, name, err)
		return err
	}
	klog.Infof("delete %s/%s cr success", namespace, name)
	return nil
}

func (r *WorkloadReconciler) IsTopOwnerReference(ownerRef []metav1.OwnerReference) bool {
	if len(ownerRef) == 0 {
		return true
	}
	for _, ref := range ownerRef {
		if *ref.Controller {
			return false
		}
	}
	return true
}

func (r *WorkloadReconciler) existCR(namespace, name string) bool {
	_, err := r.lister.Recommendations(namespace).Get(name)
	return err == nil
}

func (r *WorkloadReconciler) SetupWithManager(mgr ctrl.Manager) error {
	createAndDelete := predicate.Funcs{
		CreateFunc: func(e event.CreateEvent) bool {
			return true
		},
		DeleteFunc: func(e event.DeleteEvent) bool {
			return true
		},
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1.StatefulSet{}, builder.WithPredicates(createAndDelete)).
		Watches(&appsv1.Deployment{}, &handler.EnqueueRequestForObject{}, builder.WithPredicates(createAndDelete)).
		Watches(&appsv1.ReplicaSet{}, &handler.EnqueueRequestForObject{}, builder.WithPredicates(createAndDelete)).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).Complete(r)
}
