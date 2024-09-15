/*
Copyright 2024 by infrapot

This program is a free software product. You can redistribute it and/or
modify it under the terms of the GNU Affero General Public License (AGPL)
version 3 as published by the Free Software Foundation.

For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
*/
// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	"context"
	time "time"

	hykubev1alpha1 "github.com/infrapot/hykube/pkg/apis/hykube/v1alpha1"
	versioned "github.com/infrapot/hykube/pkg/generated/clientset/versioned"
	internalinterfaces "github.com/infrapot/hykube/pkg/generated/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/infrapot/hykube/pkg/generated/listers/hykube/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// ProviderInformer provides access to a shared informer and lister for
// Providers.
type ProviderInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.ProviderLister
}

type providerInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewProviderInformer constructs a new informer for Provider type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewProviderInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredProviderInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredProviderInformer constructs a new informer for Provider type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredProviderInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.HykubeV1alpha1().Providers(namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.HykubeV1alpha1().Providers(namespace).Watch(context.TODO(), options)
			},
		},
		&hykubev1alpha1.Provider{},
		resyncPeriod,
		indexers,
	)
}

func (f *providerInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredProviderInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *providerInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&hykubev1alpha1.Provider{}, f.defaultInformer)
}

func (f *providerInformer) Lister() v1alpha1.ProviderLister {
	return v1alpha1.NewProviderLister(f.Informer().GetIndexer())
}
