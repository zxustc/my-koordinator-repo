package frontend

import (
	"context"
	"encoding/json"
	"sort"
	"time"

	analysisv1alpha1 "github.com/koordinator-sh/koordinator/apis/analysis/v1alpha1"
	"github.com/koordinator-sh/koordinator/pkg/client/clientset/versioned/typed/analysis/v1alpha1"
	"github.com/koordinator-sh/koordinator/pkg/prediction/manager"
	"github.com/koordinator-sh/koordinator/pkg/prediction/manager/protocol"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
)

var _ StatusFetcher = &statusFetcher{}

type statusFetcher struct {
	client     v1alpha1.RecommendationsGetter
	predictMgr *manager.PredictionMgrImpl
	ctx        context.Context
	gvk        *schema.GroupVersionKind
}

type StatusFetcher interface {
	Run()
	Started() bool
	UpdateStatus()
}

func InitStatusFetcher(client v1alpha1.RecommendationsGetter, ctx context.Context, predictMgr *manager.PredictionMgrImpl) *statusFetcher {
	return &statusFetcher{
		client:     client,
		ctx:        ctx,
		predictMgr: predictMgr,
		gvk: &schema.GroupVersionKind{
			Group:   "analysis.koordinator.sh",
			Version: "v1alpha1",
			Kind:    "",
		},
	}
}

func getNextUpdate(now time.Time) time.Time {
	next := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
	if next.Before(now) {
		next = next.Add(24 * time.Hour)
	}
	return next
}

func (s *statusFetcher) Run() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(time.Until(getNextUpdate(time.Now()))):
			s.UpdateStatus()
		}
	}
}

func (s *statusFetcher) Started() bool {
	return false
}

func (s *statusFetcher) UpdateStatus() {
	listMap := s.predictMgr.GetList()
	for key, profileKey := range listMap {
		namespace, name := protocol.ParsePredictListKey(key)
		status := s.predictMgr.GetResult(profileKey)
		recommendationStatus := MapToListOfRecommendedContainerResources(status)
		recommendCrd, err := s.patchRecommendation(name, namespace, recommendationStatus)
		if err != nil {
			klog.Error("Error while patch Recommend status: %+v", err)
			return
		}
		s.predictMgr.UpdateCrdCache(profileKey, recommendCrd)
	}
}

type patchRecord struct {
	Op    string      `json:"op,inline"`
	Path  string      `json:"path,inline"`
	Value interface{} `json:"value"`
}

func (s *statusFetcher) patchRecommendation(name, namespace string, status *analysisv1alpha1.RecommendedPodStatus) (result *analysisv1alpha1.Recommendation, err error) {
	patches := []patchRecord{{
		Op:    "add",
		Path:  "/status",
		Value: status,
	}}
	bytes, err := json.Marshal(patches)
	if err != nil {
		klog.Error("Can not marshal Recommend status patches %+v, Reason: %+v", patches, err)
		return nil, err
	}
	client := s.client.Recommendations(namespace)
	return client.Patch(context.TODO(), name, types.JSONPatchType, bytes, metav1.PatchOptions{})
}

// MapToListOfRecommendedContainerResources converts the map of RecommendedContainerResources into a stable sorted list
// This can be used to get a stable sequence while ranging on the data
func MapToListOfRecommendedContainerResources(resources protocol.PredictPodStatus) *analysisv1alpha1.RecommendedPodStatus {
	containerRecommend := make([]analysisv1alpha1.RecommendedContainerStatus, 0, len(resources))
	containerNames := make([]string, 0, len(resources))

	for containerName := range resources {
		containerNames = append(containerNames, containerName)
	}

	sort.Strings(containerNames)
	for _, name := range containerNames {
		containerRecommend = append(containerRecommend, analysisv1alpha1.RecommendedContainerStatus{
			ContainerName: name,
			Resources:     resources[name].Resources,
		})
	}
	recommended := &analysisv1alpha1.RecommendedPodStatus{
		ContainerStatuses: containerRecommend,
	}
	return recommended
}

//used to update crd status
