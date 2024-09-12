/*
Copyright 2017 The Kubernetes Authors.

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

package provider

import (
	"context"
	"hykube.io/apiserver/pkg/apis/hykube"
	"hykube.io/apiserver/pkg/registry"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
)

// NewREST returns a RESTStorage object that will work against API services.
func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*registry.REST, error) {
	strategy := NewStrategy(scheme)

	store := &genericregistry.Store{
		NewFunc:                   func() runtime.Object { return &hykube.Provider{} },
		NewListFunc:               func() runtime.Object { return &hykube.ProviderList{} },
		PredicateFunc:             MatchProvider,
		DefaultQualifiedResource:  hykube.Resource("providers"),
		SingularQualifiedResource: hykube.Resource("provider"),

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,

		// TODO: define table converter that exposes more than name/creation timestamp
		TableConvertor: rest.NewDefaultTableConvertor(hykube.Resource("providers")),
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	watch, err := store.Watch(context.TODO(), &internalversion.ListOptions{})
	if err != nil {
		return nil, err
	}
	providerWatcher := watcher{
		watch:             watch,
		hashicorpProvider: NewHashicorpProvider(store),
	}
	providerWatcher.Start()

	return &registry.REST{Store: store}, nil
}
