/*
Copyright 2022 The Koordinator Authors.

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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	autoscalerv1 "k8s.io/autoscaler/vertical-pod-autoscaler/pkg/apis/autoscaling.k8s.io/v1"
	"k8s.io/kubernetes/pkg/apis/autoscaling"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// MetricPredictionSpec defines the desired state of MetricPrediction
type MetricPredictionSpec struct {

	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Metrics []MetricSource `json:"metrics,omitempty"`
}

type MetricSourceType string

const (
	MetricServerSourceType MetricSourceType = "MetricServer"
	PrometheusSourceType MetricSourceType = "Prometheus"
	// KoordinatorBuiltinSourceType MetricSourceType = "Koordinator"
)

type MetricSource struct {
	Type MetricSourceType `json:"type,omitempty"`
	APIService *MetricServerSource       `json:"apiService,omitempty"`
	Prometheus *PrometheusSource         `json:"prometheus,omitempty"`
	// Builtin    *KoordinatorBuiltinSource `json:"builtin,omitempty"`
}

type PrometheusSource struct{
	MetricName string
	GroupedLabels map[string]string
}

type ObjectReferenceType string

const (
	WorkloadObjectReferenceType ObjectReferenceType = "Workload"
	PodLabelsObjectReferenceType ObjectReferenceType = "PodLabels"
	PrometheusLabelsObjectReferenceType ObjectReferenceType = "PrometheusLabels"
)

type WorkloadObjectReference struct {
	*autoscaling.CrossVersionObjectReference `json:",inline"`
}

type PodLabelObjectReference struct {
	Selector *metav1.LabelSelector
}

type PrometheusLabelObjectReference struct {

}

type ObjectOwnerReference struct {
	Type ObjectReferenceType
	Workload *WorkloadObjectReference
	PodLabel *PodLabelObjectReference
	PrometheusLabel *PrometheusLabelObjectReference
}

type TargetObjectType string

const (
	TargetObjectContainer TargetObjectType = "Container"
	TargetObjectPrometheusLabel TargetObjectType = "PrometheusLabel"
)

type ContainerObject struct {
	Name string
	Owner ObjectOwnerReference
}

type MetricServerSource struct{
	ResourceMetricName string
	ObjectType MetricServerObjectType
	Container *ContainerObject
}

type TargetObject struct {

}

type Estimator struct {
	Distribution *DistributionEstimator `json:"distribution,omitempty"`
}

type DistributionEstimator struct {
	Quantiles []string          `json:"quantiles,omitempty"`
	Histogram *HistogramOptions `json:"histogram,omitempty"`
}

type HistogramOptions struct{

}

// MetricPredictionStatus defines the observed state of MetricPrediction
type MetricPredictionStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []EstimateCondition `json:"conditions,omitempty"`
}

type EstimateValue struct {
	Distribution *DistributionValues `json:"percentile,omitempty"`
}

type DistributionValues struct {
	Mean              resource.Quantity                `json:"mean,omitempty"`
	Quantiles         map[string]resource.Quantity     `json:"quantiles,omitempty"`
	StdDev            resource.Quantity                `json:"stddev,omitempty"`
	FirstSampleStart  metav1.Time                      `json:"firstSampleStart,omitempty"`
	LastSampleStart   metav1.Time                      `json:"lastSampleStart,omitempty"`
	TotalSamplesCount int                              `json:"totalSamplesCount,omitempty"`
	Histogram         autoscalerv1.HistogramCheckpoint `json:"histogram,omitempty"`
}

type EstimateCondition struct{}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MetricPrediction is the Schema for the metricpredictions API
// +k8s:openapi-gen=true
type MetricPrediction struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MetricPredictionSpec   `json:"spec,omitempty"`
	Status MetricPredictionStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// MetricPredictionList contains a list of MetricPrediction
type MetricPredictionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []MetricPrediction `json:"items"`
}

func init() {
	SchemeBuilder.Register(&MetricPrediction{}, &MetricPredictionList{})
}
