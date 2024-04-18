package frontend

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	analysisv1alpha1 "github.com/koordinator-sh/koordinator/apis/analysis/v1alpha1"
	"github.com/koordinator-sh/koordinator/pkg/client/clientset/versioned/typed/analysis/v1alpha1"
	"github.com/koordinator-sh/koordinator/pkg/prediction/manager"
	"github.com/koordinator-sh/koordinator/pkg/prediction/manager/apis"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/klog/v2"
)

var _ StatusFetcher = &statusFetcher{}

type statusFetcher struct {
	client       v1alpha1.RecommendationsGetter
	predictMgr   manager.PredictionManager
	ctx          context.Context
	updatePeriod int
}

type StatusFetcher interface {
	Run()
	Started() bool
	UpdateStatus()
}

func Wrapper(input string, wrapper string) string {
	return input + " " + wrapper
}

func InitStatusFetcher(client v1alpha1.RecommendationsGetter, ctx context.Context, predictMgr manager.PredictionManager) *statusFetcher {
	return &statusFetcher{
		client:     client,
		ctx:        ctx,
		predictMgr: predictMgr,
	}
}

func getNextUpdate(now time.Time) time.Time {
	next := time.Date(now.Year(), now.Month(), now.Day(), 2, 0, 0, 0, now.Location())
	if next.Before(now) {
		next = next.Add(24 * time.Hour)
	}
	return next
}

func (s *statusFetcher) metricsHandler(w http.ResponseWriter, r *http.Request) {
	promhttp.Handler().ServeHTTP(w, r)
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

func (s *statusFetcher) AsStatus(r *apis.DistributionProfilerResult) (*analysisv1alpha1.RecommendationStatus, error) {
	return nil, fmt.Errorf("has no implement")

}
func (s *statusFetcher) UpdateStatus() {
	var wg sync.WaitGroup
	sem := make(chan struct{}, 5)
	s.predictMgr.ShareState.ProfilerList.Range(func(k, v interface{}) bool {
		wg.Add(1)
		sem <- struct{}{}
		go func(k interface{}) {
			defer wg.Done()
			defer func() { <-sem }()
			key := k.(apis.ProfileKey)
			result := &apis.DistributionProfilerResult{
				DName:      key.Name(),
				DNamespace: key.Namespace(),
			}
			err := s.predictMgr.GetResult(key, result, nil)
			if err != nil {
				klog.Error("Error while get result of Profilekey: ", key, err)
			}
			status, err := s.AsStatus(result)
			if err != nil {
				klog.Error("Error while transfrom status:", err)
				return
			}
			_, err = s.patchRecommendation(key.Name(), key.Namespace(), status)
			if err != nil {
				klog.Error("Error while patch Recommend status:", err)
				return
			}
			klog.Infof("Update %s/%s status", key.Namespace(), key.Name())
		}(k)
		return true
	})
	wg.Wait()
}

type patchRecord struct {
	Op    string      `json:"op,inline"`
	Path  string      `json:"path,inline"`
	Value interface{} `json:"value"`
}

func (s *statusFetcher) patchRecommendation(name string, namespace string,
	status *analysisv1alpha1.RecommendationStatus) (result *analysisv1alpha1.Recommendation, err error) {
	patches := []patchRecord{{
		Op:    "add",
		Path:  "/status",
		Value: status,
	}}
	bytes, err := json.Marshal(patches)
	if err != nil {
		klog.Errorf("Can not marshal Recommend status patches %+v, Reason: %+v", patches, err)
		return nil, err
	}
	client := s.client.Recommendations(namespace)
	return client.Patch(context.TODO(), name, types.JSONPatchType, bytes, metav1.PatchOptions{})
}
