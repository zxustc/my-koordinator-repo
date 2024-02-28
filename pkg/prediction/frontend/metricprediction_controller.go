package frontend

import (
	"context"

	analysisv1alpha1 "github.com/koordinator-sh/koordinator/apis/analysis/v1alpha1"
	sharestate "github.com/koordinator-sh/koordinator/pkg/prediction/sharestat"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	"k8s.io/cri-api/pkg/errors"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const Name = "prediction"

// MetricPrediction reconciles a MetricPrediction object
type MetricPredictionReconcile struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=analysis.koordinator.sh,resources=metricpredictions,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=analysis.koordinator.sh,resources=metricpredictions/status,verbs=get;update;patch
func (r *MetricPredictionReconcile) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var metricPrediction analysisv1alpha1.MetricPrediction

	_ = log.FromContext(ctx, "metric-prediction-reconcile", req.NamespacedName)
	event := req.String()
	switch event {
	case "CREATE":
		if err := r.Client.Get(ctx, req.NamespacedName, &metricPrediction); err != nil {
			// If no cr, retry.
			if !errors.IsNotFound(err) {
				klog.Errorf("failed to Get metricPrediction %v, error: %v", metricPrediction, err)
			}
			return ctrl.Result{Requeue: true}, err
		} else {
			sharestate.UpdateShareState(req.NamespacedName)
			return ctrl.Result{}, nil
		}
	case "DELETE":
		if err := r.Client.Get(ctx, req.NamespacedName, &metricPrediction); err != nil {
			if errors.IsNotFound(err) {
				sharestate.DeleteShareState(req.NamespacedName)
				return ctrl.Result{}, nil
			} else {
				return ctrl.Result{Requeue: true}, err
			}
		} else {
			return ctrl.Result{Requeue: true}, err
		}
	case "UPDATE":
		if err := r.Client.Get(ctx, req.NamespacedName, &metricPrediction); err != nil {
			// If no cr, retry.
			if !errors.IsNotFound(err) {
				klog.Errorf("failed to Get metricPrediction %v, error: %v", metricPrediction, err)
			}
			return ctrl.Result{Requeue: true}, err
		} else {
			sharestate.UpdateShareState(req.NamespacedName)
			return ctrl.Result{}, nil
		}
	default:
		return ctrl.Result{}, nil
	}
}

func Add(mgr ctrl.Manager) error {
	reconciler := MetricPredictionReconcile{
		Client:   mgr.GetClient(),
		Scheme:   mgr.GetScheme(),
		Recorder: mgr.GetEventRecorderFor("metric-prediction"),
	}
	return reconciler.SetupWithManager(mgr)
}

// SetupWithManager sets up the controller with the Manager.
func (r *MetricPredictionReconcile) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&analysisv1alpha1.MetricPrediction{}, builder.WithPredicates(predicate.GenerationChangedPredicate{})).
		Complete(r)
}
