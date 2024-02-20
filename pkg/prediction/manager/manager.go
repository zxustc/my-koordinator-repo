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
	Update(...protocol.PredictionProfileKey) error
	GetResult(key protocol.PredictionProfileKey) (protocol.PredictionResult, error)
}

var _ PredictionManager = &predictionMgrImpl{}

type predictionMgrImpl struct {
	metricsRepo     metricscollector.MetricsRepository
	checkpoint      checkpoint.Checkpoint
	workloadFetcher workloadfetcher.WorkloadFetcher
	profiler        profiler.Profiler
}

func (p *predictionMgrImpl) Run() error {
	// run checkpoint to load all history data
	// start workload fetcher waiting for workloads
	// start metrics repo ready for collect
	// start profiler to calculate each model
	panic("implement me")
}

func (p *predictionMgrImpl) Started() bool {
	// return true only if all components are started
	panic("implement me")
}

func (p *predictionMgrImpl) Register(profiles ...protocol.PredictionProfile) error {
	// return error if not started
	// add workload to WorkloadFetcher
	// subscribe metric to metricsRepo
	// add model to Profiler
	panic("implement me")
}

func (p *predictionMgrImpl) Unregister(keys ...protocol.PredictionProfileKey) error {
	// return error if not started
	// remove workload from WorkloadFetcher
	// unsubscribe metric from metricsRepo
	// remove model from Profiler
	panic("implement me")
}

func (p *predictionMgrImpl) GetResult(key protocol.PredictionProfileKey) (protocol.PredictionResult, error) {
	//TODO implement me
	panic("implement me")
}
