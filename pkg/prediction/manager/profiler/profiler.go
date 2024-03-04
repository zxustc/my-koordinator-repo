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

package profiler

import "github.com/koordinator-sh/koordinator/pkg/prediction/manager/protocol"

type Profiler interface {
	Run() error
	Started() bool
	GetResult(key protocol.PredictionProfileKey) (protocol.PredictionResult, error)
}

var _ Profiler = &profilerImpl{}

type profilerImpl struct {
	models map[ModelKey]Model
}

func InitProfiler() *profilerImpl {
	return &profilerImpl{}
}

func (p *profilerImpl) Run() error {
	// get groupings (pod/container belongs to same workload) list from workload fetcher
	// create/update(args) model for each grouping
	// load history from checkpoint if exist for new model
	// for each model, get metric from metric repo and feed samples to model
	// save checkpoint for each model
	panic("implement me")
}

func (p *profilerImpl) Started() bool {
	//TODO implement me
	panic("implement me")
}

func (p *profilerImpl) GetResult(key protocol.PredictionProfileKey) (protocol.PredictionResult, error) {
	//TODO implement me
	panic("implement me")
}
