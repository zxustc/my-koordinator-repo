package protocol

import (
	analysisv1 "github.com/koordinator-sh/koordinator/apis/analysis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
)

func (cluster *ClusterState) getLabelSetKey(labelSet labels.Set) labelSetKey {
	labelSetKey := labelSetKey(labelSet.String())
	cluster.labelSetMap[labelSetKey] = labelSet
	return labelSetKey
}

// PodID contains information needed to identify a Pod within a cluster.
type PodID struct {
	// Namespaces where the Pod is defined.
	Namespace string
	// PodName is the name of the pod unique within a namespace.
	PodName string
}

// ContainerID contains information needed to identify a Container within a cluster.
type ContainerID struct {
	PodID
	// ContainerName is the name of the container, unique within a pod.
	ContainerName string
}

// Analysis crd information needed to identify a metric prediction API object within a cluster.
type AnalysisTargetID struct {
	Namespace            string
	MetricPredictionName string
}

type labelSetKey string

// Map of label sets keyed by their string representation.
type labelSetMap map[labelSetKey]labels.Set

type PodState struct {
	// Unique id of the Pod.
	ID PodID
	// Set of labels attached to the Pod.
	labelSetKey labelSetKey
	// Containers that belong to the Pod, keyed by the container name.
	Containers map[string]*ContainerState
	// PodPhase describing current life cycle phase of the Pod.
	Phase corev1.PodPhase
}

// ClusterState is used to cache workling lists for workers.
type ClusterState struct {
	// workloads Need to Profile
	workLists []predictionProfileKey

	// Pods in the cluster.
	Pods map[PodID]*PodState
	// AnalysisTarget CRD objects in the cluster.
	analysisTarget map[AnalysisTargetID]*analysisv1.AnalysisTarget
	// All container aggregations where the usage samples are stored.
	aggregateInfo map[PredictionProfileKey]PredictionProfileState
	// Map with all label sets used by the aggregations. It serves as a cache
	// that allows to quickly access labels.Set corresponding to a labelSetKey.
	labelSetMap labelSetMap
}

type ContainerState struct{}
