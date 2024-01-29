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

import v1 "k8s.io/api/core/v1"

// MetricSourceType defines the type of metric source
type MetricSourceType string

const (
	// MetricSourceTypeMetricServer means the metric is from metric server
	MetricSourceTypeMetricServer MetricSourceType = "metricServer"
	// MetricSourceTypePrometheus means the metric is from prometheus
	MetricSourceTypePrometheus MetricSourceType = "prometheus"
)

// MetricSpec defines the metric to be analysis
type MetricSpec struct {
	// Source defines the source of metric, which can be metric server or prometheus
	Source MetricSourceType `json:"source"`
	// MetricServer defines the metric server source, which is effective when source is metric server
	MetricServer *MetricServerSource `json:"metricServer,omitempty"`
	// Prometheus defines the prometheus source, which is effective when source is prometheus
	Prometheus *PrometheusMetricSource `json:"prometheus,omitempty"`
}

// MetricServerSource defines the metric server source
type MetricServerSource struct {
	// Resources indicates the resources to be analyzed, only cpu and memory supported for metric server
	Resources []v1.ResourceName `json:"resources,omitempty"`
}

// PrometheusMetricSource defines the prometheus metric source
type PrometheusMetricSource struct {
	// Metrics defines the prometheus metrics to be analyzed
	Metrics []PrometheusMetric `json:"metrics,omitempty"`
}

// PrometheusMetric defines the prometheus metric to be analyzed
type PrometheusMetric struct {
	// Resource defines the key of resource to be analyzed
	Resource v1.ResourceName `json:"name,omitempty"`
	// MetricName is the name of prometheus metric, such as container_cpu_usage
	MetricName string `json:"metric,omitempty"`
}
