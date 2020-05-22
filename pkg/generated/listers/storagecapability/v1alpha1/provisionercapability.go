/*
Copyright The Kubernetes Authors.

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
	v1alpha1 "github.com/kubesphere/storage-capability/pkg/apis/storagecapability/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ProvisionerCapabilityLister helps list ProvisionerCapabilities.
type ProvisionerCapabilityLister interface {
	// List lists all ProvisionerCapabilities in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.ProvisionerCapability, err error)
	// Get retrieves the ProvisionerCapability from the index for a given name.
	Get(name string) (*v1alpha1.ProvisionerCapability, error)
	ProvisionerCapabilityListerExpansion
}

// provisionerCapabilityLister implements the ProvisionerCapabilityLister interface.
type provisionerCapabilityLister struct {
	indexer cache.Indexer
}

// NewProvisionerCapabilityLister returns a new ProvisionerCapabilityLister.
func NewProvisionerCapabilityLister(indexer cache.Indexer) ProvisionerCapabilityLister {
	return &provisionerCapabilityLister{indexer: indexer}
}

// List lists all ProvisionerCapabilities in the indexer.
func (s *provisionerCapabilityLister) List(selector labels.Selector) (ret []*v1alpha1.ProvisionerCapability, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ProvisionerCapability))
	})
	return ret, err
}

// Get retrieves the ProvisionerCapability from the index for a given name.
func (s *provisionerCapabilityLister) Get(name string) (*v1alpha1.ProvisionerCapability, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("provisionercapability"), name)
	}
	return obj.(*v1alpha1.ProvisionerCapability), nil
}
