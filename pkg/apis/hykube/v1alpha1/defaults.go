/*
 * Copyright 2024 by infrapot
 *
 * This program is a free software product. You can redistribute it and/or
 * modify it under the terms of the GNU Affero General Public License (AGPL)
 * version 3 as published by the Free Software Foundation.
 *
 * For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
 */

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/runtime"
)

func addDefaultingFuncs(scheme *runtime.Scheme) error {
	return RegisterDefaults(scheme)
}

// SetDefaults_ProviderSpec sets defaults for Provider spec
func SetDefaults_ProviderSpec(obj *ProviderSpec) {

}

// SetDefaults_PlanSpec sets defaults for Plan spec
func SetDefaults_PlanSpec(p *PlanSpec) {

}
