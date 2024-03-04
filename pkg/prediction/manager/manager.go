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

	"github.com/koordinator-sh/koordinator/pkg/prediction/manager/checkpoint"
	"github.com/koordinator-sh/koordinator/pkg/prediction/manager/metricscollector"
	"github.com/koordinator-sh/koordinator/pkg/prediction/manager/profiler"
	"github.com/koordinator-sh/koordinator/pkg/prediction/manager/protocol"
	"github.com/koordinator-sh/koordinator/pkg/prediction/manager/workloadfetcher"
)

type PredictionManager interface {
	Run() error
	Started() bool
	Register(...protocol.PredictionProfile) error
	Unregister(...protocol.PredictionProfileKey) error
	GetResult(key protocol.PredictionProfileKey) (protocol.PredictionResult, error)
}

var _ PredictionManager = &PredictionMgrImpl{}

type PredictionMgrImpl struct {
	metricsRepo     metricscollector.MetricsRepository
	checkpoint      checkpoint.Checkpoint
	workloadFetcher workloadfetcher.WorkloadFetcher
	profiler        profiler.Profiler
}

func InitPredictMgr() *PredictionMgrImpl {

	m := metricscollector.InitMetricRepo()
	c := checkpoint.InitCheckpoint()
	w := workloadfetcher.InitWorkloadfetcher()
	p := profiler.InitProfiler()

	predictMgr := &PredictionMgrImpl{
		metricsRepo:     m,
		checkpoint:      c,
		workloadFetcher: w,
		profiler:        p,
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

func (p *PredictionMgrImpl) Register(profiles ...protocol.PredictionProfile) error {
	// return error if not started
	if !p.Started() {
		return fmt.Errorf("Profiler did not start.")
	}
	// add workload to WorkloadFetcher
	p.workloadFetcher.AddWorkloads(profiles)
	// subscribe metric to metricsRepo
	p.metricsRepo.Register()
	// add model to Profiler
	panic("implement me")
}

func (p *PredictionMgrImpl) Unregister(keys ...protocol.PredictionProfileKey) error {
	// return error if not started
	// remove workload from WorkloadFetcher
	// unsubscribe metric from metricsRepo
	// remove model from Profiler
	p.workloadFetcher.RemoveWorkloads()
	p.metricsRepo.Unregister()
	panic("implement me")
}

func (p *PredictionMgrImpl) GetResult(key protocol.PredictionProfileKey) (protocol.PredictionResult, error) {
	//get result and fill status
	//TODO implement me
	res, err := p.profiler.GetResult(key)
	if err != nil {
		return nil, err
	}
	return res, err
}
