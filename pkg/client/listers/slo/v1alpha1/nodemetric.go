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
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/koordinator-sh/koordinator/apis/slo/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// NodeMetricLister helps list NodeMetrics.
// All objects returned here must be treated as read-only.
type NodeMetricLister interface {
	// List lists all NodeMetrics in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.NodeMetric, err error)
	// Get retrieves the NodeMetric from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.NodeMetric, error)
	NodeMetricListerExpansion
}

// nodeMetricLister implements the NodeMetricLister interface.
type nodeMetricLister struct {
	indexer cache.Indexer
}

// NewNodeMetricLister returns a new NodeMetricLister.
func NewNodeMetricLister(indexer cache.Indexer) NodeMetricLister {
	return &nodeMetricLister{indexer: indexer}
}

// List lists all NodeMetrics in the indexer.
func (s *nodeMetricLister) List(selector labels.Selector) (ret []*v1alpha1.NodeMetric, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.NodeMetric))
	})
	return ret, err
}

// Get retrieves the NodeMetric from the index for a given name.
func (s *nodeMetricLister) Get(name string) (*v1alpha1.NodeMetric, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("nodemetric"), name)
	}
	return obj.(*v1alpha1.NodeMetric), nil
}