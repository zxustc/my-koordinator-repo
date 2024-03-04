package protocol

import (
	"sync"

	analysisv1 "github.com/koordinator-sh/koordinator/apis/analysis/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
)

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

type ContainerState struct{}

type ClusterState struct {
	listMutex          sync.Mutex
	crdMutex           sync.Mutex
	resultMutex        sync.Mutex
	RecommendationList PredictListMap
	Result             PredictPodStatusMap
	RecommendCrds      RecommendCrds //crd cache?
}

func NewClusterState() *ClusterState {
	clusterState := &ClusterState{
		RecommendationList: make(PredictListMap),
		Result:             make(PredictPodStatusMap),
		RecommendCrds:      make(RecommendCrds),
	}
	return clusterState
}

func (c *ClusterState) AddToList(rec analysisv1.RecommendationSpec, id types.NamespacedName) {
	defer c.listMutex.Unlock()
	key := MakePredictListKey(id)
	profileKey := InitPredictionProfileKey(rec, id)
	c.listMutex.Lock()
	if _, ok := c.RecommendationList[key]; !ok {
		c.RecommendationList[key] = profileKey
	}
}

func (c *ClusterState) RemoveFromList(id types.NamespacedName) {
	defer c.listMutex.Unlock()
	listID := MakePredictListKey(id)
	c.listMutex.Lock()
	delete(c.RecommendationList, listID)
}

func (c *ClusterState) UpdateCrdCache(key PredictionProfileKey, cache *analysisv1.Recommendation) {
	defer c.crdMutex.Unlock()
	c.crdMutex.Lock()
	c.RecommendCrds[key] = cache
}

func (c *ClusterState) DeleteCrdCache(key PredictionProfileKey) {
	defer c.crdMutex.Unlock()
	c.crdMutex.Lock()
	delete(c.RecommendCrds, key)
}
