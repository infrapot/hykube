/*
Copyright 2024 by infrapot

This program is a free software product. You can redistribute it and/or
modify it under the terms of the GNU Affero General Public License (AGPL)
version 3 as published by the Free Software Foundation.

For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
*/
// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/infrapot/hykube/pkg/apis/hykube/v1alpha1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/listers"
	"k8s.io/client-go/tools/cache"
)

// ProviderLister helps list Providers.
// All objects returned here must be treated as read-only.
type ProviderLister interface {
	// List lists all Providers in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Provider, err error)
	// Providers returns an object that can list and get Providers.
	Providers(namespace string) ProviderNamespaceLister
	ProviderListerExpansion
}

// providerLister implements the ProviderLister interface.
type providerLister struct {
	listers.ResourceIndexer[*v1alpha1.Provider]
}

// NewProviderLister returns a new ProviderLister.
func NewProviderLister(indexer cache.Indexer) ProviderLister {
	return &providerLister{listers.New[*v1alpha1.Provider](indexer, v1alpha1.Resource("provider"))}
}

// Providers returns an object that can list and get Providers.
func (s *providerLister) Providers(namespace string) ProviderNamespaceLister {
	return providerNamespaceLister{listers.NewNamespaced[*v1alpha1.Provider](s.ResourceIndexer, namespace)}
}

// ProviderNamespaceLister helps list and get Providers.
// All objects returned here must be treated as read-only.
type ProviderNamespaceLister interface {
	// List lists all Providers in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Provider, err error)
	// Get retrieves the Provider from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Provider, error)
	ProviderNamespaceListerExpansion
}

// providerNamespaceLister implements the ProviderNamespaceLister
// interface.
type providerNamespaceLister struct {
	listers.ResourceIndexer[*v1alpha1.Provider]
}
