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

// ReferenceType defines the type of object reference.
type ReferenceType string

const (
	// ProviderReferenceType is used for Provider references.
	ProviderReferenceType = ReferenceType("Provider")
)

// ProviderSpec is the specification of a Provider.
type ProviderSpec struct {
	// A name of another provider. TODO: is it needed??
	ProviderReference string
	// The reference type.
	ReferenceType ReferenceType

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
