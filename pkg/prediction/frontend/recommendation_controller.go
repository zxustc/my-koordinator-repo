/*
 * register reconciler of recommendation crd.
 */

package frontend

import (
	"context"

	analysisv1alpha1 "github.com/koordinator-sh/koordinator/apis/analysis/v1alpha1"
	"github.com/koordinator-sh/koordinator/pkg/prediction/manager"
	"github.com/koordinator-sh/koordinator/pkg/prediction/manager/apis"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

// Recommendation reconciles a Recommendation object
type RecommendationReconcile struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
	Manager  manager.PredictionManager
}

// +kubebuilder:rbac:groups=analysis.koordinator.sh,resources=metricpredictions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=analysis.koordinator.sh,resources=metricpredictions/status,verbs=get;update;patch
func (r *RecommendationReconcile) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var recommendation analysisv1alpha1.Recommendation
	_ = log.FromContext(ctx, "recommendation-reconcile", req.NamespacedName)
	if err := r.Client.Get(ctx, req.NamespacedName, &recommendation); err != nil {
		profileKey := apis.MakeProfileKey(req.NamespacedName)
		// unreg key
		err := r.Manager.Unregister(profileKey)
		if err != nil {
			return ctrl.Result{Requeue: true}, err
		}
		return ctrl.Result{}, nil
	}
	profileKey := apis.MakeProfileKey(req.NamespacedName)
	profiler := apis.MakeProfilerSpec(profileKey, recommendation)
	err := r.Manager.Register(profileKey, profiler)
	if err != nil {
		return ctrl.Result{Requeue: true}, err
	}

	return ctrl.Result{}, nil
}

func Add(mgr ctrl.Manager, predictionManager manager.PredictionManager) error {
	reconciler := &RecommendationReconcile{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("metric-prediction"),
		Manager:  predictionManager,
	}
	return reconciler.SetupWithManager(mgr)
}

// SetupWithManager sets up the controller with the Manager.
func (r *RecommendationReconcile) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&analysisv1alpha1.Recommendation{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}
