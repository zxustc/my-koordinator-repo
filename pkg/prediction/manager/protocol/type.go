package protocol

import (
	"strings"

	analysisv1 "github.com/koordinator-sh/koordinator/apis/analysis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ PredictionProfileKey = &predictionProfileKey{}

// predictionProfileKey is used to collect infos to position a unique profile task
type predictionProfileKey struct {
	TargetType analysisv1.RecommendationTargetType
	//Pod selector that determines profile pods, used with type "workload".
	Workload *analysisv1.CrossVersionObjectReference
	//Pod selector that determines profile pods, used with type "podselector".
	PodSelector *metav1.LabelSelector
}

// ProfileID
type RecommendationID struct {
	//metricPredication CR’s namespace
	Namespace string
	//metricPredication CR’s name
	Name string
}

func (p predictionProfileKey) Type() string {
	return string(p.TargetType)
}

func (p predictionProfileKey) WorkloadName() string {
	if p.TargetType == "workload" {
		return p.Workload.Kind
	} else {
		return ""
	}
}

func (p predictionProfileKey) WorkloadKind() string {
	if p.TargetType == "workload" {
		return p.Workload.Name
	} else {
		return ""
	}
}

func (p predictionProfileKey) LabelSelector() *metav1.LabelSelector {
	if p.TargetType == "workload" {
		return nil
	} else {
		return p.PodSelector
	}
}

func InitPredictionProfileKey(recommendation analysisv1.RecommendationSpec, id types.NamespacedName) PredictionProfileKey {
	return predictionProfileKey{
		TargetType:  recommendation.Target.Type,
		Workload:    recommendation.Target.Workload,
		PodSelector: recommendation.Target.PodSelector,
	}
}

type PredictContainerResource struct {
	Resources corev1.ResourceList
}

type PredictListMap map[string]PredictionProfileKey
type PredictPodStatus map[string]PredictContainerResource
type PredictPodStatusMap map[PredictionProfileKey]PredictPodStatus
type RecommendCrds map[PredictionProfileKey]*analysisv1.Recommendation

func MakePredictListKey(key types.NamespacedName) string {
	value := key.Namespace + "-" + key.Name
	return value
}

func ParsePredictListKey(key string) (string, string) {
	parts := strings.Split(key, "-")
	if len(parts) < 1 {
		return "", ""
	}
	if len(parts) == 1 {
		return parts[0], ""
	}
	return parts[0], parts[1]
}
