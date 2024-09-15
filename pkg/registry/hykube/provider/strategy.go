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
	"fmt"
	"github.com/infrapot/hykube/pkg/apis/hykube"
	"github.com/infrapot/hykube/pkg/apis/hykube/validation"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
)

// NewStrategy creates and returns a providerStrategy instance
func NewStrategy(typer runtime.ObjectTyper) providerStrategy {
	return providerStrategy{
		typer,
		names.SimpleNameGenerator,
	}
}

// GetAttrs returns labels.Set, fields.Set, and error in case the given runtime.Object is not a Provider
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	apiServer, ok := obj.(*hykube.Provider)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Provider")
	}
	return labels.Set(apiServer.ObjectMeta.Labels), SelectableFields(apiServer), nil
}

// MatchProvider is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiServer only interested in specific labels/fields.
func MatchProvider(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *hykube.Provider) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

type providerStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (providerStrategy) NamespaceScoped() bool {
	return true
}

func (p providerStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	obj.(*hykube.Provider).Status = "Initiated"
}

func (providerStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (providerStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	provider := obj.(*hykube.Provider)
	return validation.ValidateProvider(provider)
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (providerStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (providerStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (providerStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (providerStrategy) Canonicalize(obj runtime.Object) {
}

func (providerStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnUpdate returns warnings for the given update.
func (providerStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}
