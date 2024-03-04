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

package protocol

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

type PredictionProfileKey interface {
	//return type "workload" or "podSelector"
	Type() string
	//Stage Name
	StageName() string
	//WorkloadName, Only valid in "workload"
	WorkloadName() string
	//WorkloadKind, Only valid in "workload"
	WorkloadKind() string
	//Label Selector
	LabelSelector() *metav1.LabelSelector
	//Metric Source Name
	MetricSourceName() string
	//Profiler Name
	ProfilerName() []string
}

type PredictionProfile interface {
	PredictionProfileKey
}

type PredictionResult interface {
	PredictionProfileKey
}
