/*
 Copyright 2024 The Koordinator Authors.

 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package manager

import (
	"fmt"

	analysisv1alpha1 "github.com/koordinator-sh/koordinator/apis/analysis/v1alpha1"

	"github.com/koordinator-sh/koordinator/pkg/prediction/manager/checkpoint"
	"github.com/koordinator-sh/koordinator/pkg/prediction/manager/metricscollector"
	"github.com/koordinator-sh/koordinator/pkg/prediction/manager/profiler"
	"github.com/koordinator-sh/koordinator/pkg/prediction/manager/protocol"
	"github.com/koordinator-sh/koordinator/pkg/prediction/manager/workloadfetcher"
	"k8s.io/apimachinery/pkg/types"
)

type PredictionManager interface {
	Run() error
	Started() bool
	Register(analysisv1alpha1.RecommendationSpec, types.NamespacedName) error
	Unregister(types.NamespacedName) error
	GetResult(protocol.PredictionProfileKey) protocol.PredictPodStatus
}

var _ PredictionManager = &PredictionMgrImpl{}

// cluster state/ share state
type PredictionMgrImpl struct {
	metricsRepo     metricscollector.MetricsRepository
	checkpoint      checkpoint.Checkpoint
	workloadFetcher workloadfetcher.WorkloadFetcher
	profiler        profiler.Profiler
	clusterState    *protocol.ClusterState
}

func InitPredictMgr() *PredictionMgrImpl {
	predictMgr := &PredictionMgrImpl{
		metricsRepo:     metricscollector.InitMetricRepo(),
		checkpoint:      checkpoint.InitCheckpoint(),
		workloadFetcher: workloadfetcher.InitWorkloadfetcher(),
		profiler:        profiler.InitProfiler(),
		clusterState:    protocol.NewClusterState(),
	}
	return predictMgr
}

func (p *PredictionMgrImpl) Run() error {
	// run checkpoint to load all history data
	// start workload fetcher waiting for workloads
	// start metrics repo ready for collect
	// start profiler to calculate each model
	panic("implement me")
}

func (p *PredictionMgrImpl) Started() bool {
	// return true only if all components are started
	panic("implement me")
}

func (p *PredictionMgrImpl) Register(rec analysisv1alpha1.RecommendationSpec, id types.NamespacedName) error {
	// return error if not started
	if !p.Started() {
		return fmt.Errorf("Profiler did not start.")
	}
	p.clusterState.AddToList(rec, id)
	// subscribe metric to metricsRepo
	p.metricsRepo.Register()
	// add model to Profiler
	panic("implement me")
}

func (p *PredictionMgrImpl) Unregister(id types.NamespacedName) error {
	// return error if not started
	// remove workload from WorkloadFetcher
	// unsubscribe metric from metricsRepo
	// remove model from Profiler
	p.clusterState.RemoveFromList(id)
	p.workloadFetcher.RemoveWorkloads()
	p.metricsRepo.Unregister()
	panic("implement me")
}

func (p *PredictionMgrImpl) GetResult(key protocol.PredictionProfileKey) protocol.PredictPodStatus {
	return p.clusterState.Result[key]
}

func (p *PredictionMgrImpl) AddtoList(id types.NamespacedName, key protocol.PredictionProfileKey) {
}

func (p *PredictionMgrImpl) GetList() protocol.PredictListMap {
	return p.clusterState.RecommendationList
}

func (p *PredictionMgrImpl) UpdateCrdCache(key protocol.PredictionProfileKey, cache *analysisv1alpha1.Recommendation) {
	p.clusterState.UpdateCrdCache(key, cache)
}
