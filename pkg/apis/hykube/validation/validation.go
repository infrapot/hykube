/*
Copyright 2016 The Kubernetes Authors.

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

package validation

import (
	"hykube.io/apiserver/pkg/apis/hykube"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateProvider validates a Provider.
func ValidateProvider(f *hykube.Provider) field.ErrorList {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, ValidateProviderSpec(&f.Spec, field.NewPath("spec"))...)

	return allErrs
}

// ValidateProviderSpec validates a ProviderSpec.
func ValidateProviderSpec(s *hykube.ProviderSpec, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if len(s.ProviderReference) == 0 && s.ReferenceType == hykube.ProviderReferenceType {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("providerReference"), s.ProviderReference, "cannot be empty if referenceType is Provider"))
	}

	if len(s.ReferenceType) != 0 && s.ReferenceType != hykube.ProviderReferenceType {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("referenceType"), s.ReferenceType, "must be Provider"))
	}

	return allErrs
}