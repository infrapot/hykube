/*
Copyright 2024 by infrapot

This program is a free software product. You can redistribute it and/or
modify it under the terms of the GNU Affero General Public License (AGPL)
version 3 as published by the Free Software Foundation.

For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
*/

package hykube

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ProviderList is a list of Provider objects.
type ProviderList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Provider
}

// ProviderSpec is the specification of a Provider.
type ProviderSpec struct {
	DownloadName string
	Version      *string
	DownloadUrl  *string
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Provider struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec     ProviderSpec
	Status   string
	Filename string
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PlanList is a list of Plan objects.
type PlanList struct {
	metav1.TypeMeta
	metav1.ListMeta

	Items []Plan
}

// PlanSpec is the specification of a Provider.
type PlanSpec struct {
	// TODO provider specific configuration
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Plan struct {
	metav1.TypeMeta
	metav1.ObjectMeta

	Spec   PlanSpec
	Status string
}
