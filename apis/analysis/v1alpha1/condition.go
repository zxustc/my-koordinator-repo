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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConditionType string

const (
	// LowConfidence indicates the low confidence for the current forecasting result.
	LowConfidence ConditionType = "LowConfidence"
	// NoPodsMatched indicates that the current description didn't match any objects.
	NoPodsMatched ConditionType = "NoObjectsMatched"
	// FetchingHistory indicates that forecaster is in the process of loading additional
	// history samples.
	FetchingHistory ConditionType = "FetchingHistory"
	// ConfigDeprecated indicates that this configuration is deprecated and will stop being
	// supported soon.
	ConfigDeprecated ConditionType = "ConfigDeprecated"
	// ConfigUnsupported indicates that this configuration is unsupported and will not be provided for it.
	ConfigUnsupported ConditionType = "ConfigUnsupported"
)

type Condition struct {
	// type describes the current condition
	Type ConditionType `json:"type"`
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
