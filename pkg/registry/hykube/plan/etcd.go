/*
 * Copyright 2024 by infrapot
 *
 * This program is a free software product. You can redistribute it and/or
 * modify it under the terms of the GNU Affero General Public License (AGPL)
 * version 3 as published by the Free Software Foundation.
 *
 * For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
 */

package plan

import (
	"context"
	"github.com/infrapot/hykube/pkg/apis/hykube"
	"github.com/infrapot/hykube/pkg/registry"
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
		NewFunc:                   func() runtime.Object { return &hykube.Plan{} },
		NewListFunc:               func() runtime.Object { return &hykube.PlanList{} },
		PredicateFunc:             MatchPlan,
		DefaultQualifiedResource:  hykube.Resource("plans"),
		SingularQualifiedResource: hykube.Resource("plan"),

		CreateStrategy: strategy,
		UpdateStrategy: strategy,
		DeleteStrategy: strategy,

		// TODO: define table converter that exposes more than name/creation timestamp
		TableConvertor: rest.NewDefaultTableConvertor(hykube.Resource("plans")),
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}

	watch, err := store.Watch(context.TODO(), &internalversion.ListOptions{})
	if err != nil {
		return nil, err
	}
	watcher := watcher{
		watch:         watch,
		hashicorpPlan: NewHashicorpPlan(store),
	}
	watcher.Start()

	return &registry.REST{Store: store}, nil
}
