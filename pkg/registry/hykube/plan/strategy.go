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

// NewStrategy creates and returns a planStrategy instance
func NewStrategy(typer runtime.ObjectTyper) planStrategy {
	return planStrategy{
		typer,
		names.SimpleNameGenerator,
	}
}

// GetAttrs returns labels.Set, fields.Set, and error in case the given runtime.Object is not the entity
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	apiServer, ok := obj.(*hykube.Plan)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Plan")
	}
	return labels.Set(apiServer.ObjectMeta.Labels), SelectableFields(apiServer), nil
}

// MatchPlan is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiServer only interested in specific labels/fields.
func MatchPlan(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *hykube.Plan) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

type planStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (planStrategy) NamespaceScoped() bool {
	return true
}

func (p planStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	obj.(*hykube.Plan).Status = "Initiated"
}

func (planStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (planStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	entity := obj.(*hykube.Plan)
	return validation.ValidatePlan(entity)
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (planStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (planStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (planStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (planStrategy) Canonicalize(obj runtime.Object) {
}

func (planStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnUpdate returns warnings for the given update.
func (planStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}
