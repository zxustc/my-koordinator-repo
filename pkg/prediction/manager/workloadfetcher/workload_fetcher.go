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

package workloadfetcher

type WorkloadFetcher interface {
	Run() error
	Started() bool
	AddWorkloads(...Workload)
	RemoveWorkloads(...Workload)
}

var _ WorkloadFetcher = &workloadFetcherImpl{}

type workloadFetcherImpl struct {
}

func InitWorkloadfetcher() *workloadFetcherImpl {
	return &workloadFetcherImpl{}
}
func (w *workloadFetcherImpl) Run() error {
	//TODO implement me
	//Start to run the workload Fetcher
	panic("implement me")
}

func (w *workloadFetcherImpl) Started() bool {
	//TODO implement me
	//If workloadfetcher is working or not?
	panic("implement me")
}

func (w *workloadFetcherImpl) AddWorkloads(workload ...Workload) {
	//TODO implement me
	//Add new workload to workfetcher working list.
	panic("implement me")
}

func (w *workloadFetcherImpl) RemoveWorkloads(workload ...Workload) {
	//TODO implement me
	panic("implement me")
	//remove workload from workfetcher working list.
}

type Workload interface {
}
