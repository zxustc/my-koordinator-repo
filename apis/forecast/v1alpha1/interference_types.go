/*
 Copyright 2023 The Koordinator Authors.
 
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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MetricSourceType string

const (
	PrometheusSourceType MetricSourceType = "Prometheus"
	MetricServerSourceType MetricSourceType = "MetricServer"
	KoordinatorBuiltinSourceType MetricSourceType = "Koordinator"
)

type MetricSource struct {
	Type MetricSourceType `json:"type,omitempty"`
}

type EstimatorType string

const (
	DistributionEstimatorType EstimatorType = "Distribution"
)

type DistributionEstimator struct {

}

type Estimator struct {
	Type EstimatorType `json:"type,omitempty"`
	Distribution *DistributionEstimator `json:"distribution,omitempty"`
}

// InterferenceDetectSpec defines the desired state of InterferenceDetect
type InterferenceDetectSpec struct {
	Metric MetricSource `json:"metric,omitempty"`
	Estimator Estimator `json:"estimator,omitempty"`
}

type DistributionValues struct {
	Mean              resource.Quantity                `json:"mean,omitempty"`
	Quantiles         map[string]resource.Quantity     `json:"quantiles,omitempty"`
	StdDev            resource.Quantity                `json:"stddev,omitempty"`
	FirstSampleStart  metav1.Time                      `json:"firstSampleStart,omitempty"`
	LastSampleStart   metav1.Time                      `json:"lastSampleStart,omitempty"`
	TotalSamplesCount int                              `json:"totalSamplesCount,omitempty"`
}

type InterferenceConditionType ConditionType

const (
	InterferenceDetected InterferenceConditionType = "InterferenceDetected"
)

type InterferenceCondition struct {
	// type describes the current condition
	Type InterferenceConditionType `json:"type"`
	// status is the status of the condition (True, False, Unknown)
	Status corev1.ConditionStatus `json:"status"`
	// lastTransitionTime is the last time the condition transitioned from one status to another
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty"`
	// reason is the reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// message is a human-readable explanation containing details about the transition
	// +optional
	Message string `json:"message,omitempty"`
}

type ContainerInterferenceDetail struct {
	MetricDistributions []DistributionValues `json:"metricDistributions,omitempty"`
}

type InterferenceDetectStatus struct {
	Conditions []InterferenceCondition `json:"conditions,omitempty"`
	Containers []ContainerInterferenceDetail `json:"containers,omitempty"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InterferenceDetect is the Schema for the metricpredictions API
// +k8s:openapi-gen=true
type InterferenceDetect struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   InterferenceDetectSpec   `json:"spec,omitempty"`
	Status InterferenceDetectStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// InterferenceDetectList contains a list of InterferenceDetect
type InterferenceDetectList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []InterferenceDetect `json:"items"`
}

func init() {
	SchemeBuilder.Register(&InterferenceDetect{}, &InterferenceDetectList{})
}
