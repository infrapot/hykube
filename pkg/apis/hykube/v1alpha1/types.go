/*
Copyright 2024 by infrapot

This program is a free software product. You can redistribute it and/or
modify it under the terms of the GNU Affero General Public License (AGPL)
version 3 as published by the Free Software Foundation.

For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
*/

package v1alpha1

import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:prerelease-lifecycle-gen:introduced=1.0
// +k8s:prerelease-lifecycle-gen:removed=1.10

// ProviderList is a list of Provider objects.
type ProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []Provider `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type ReferenceType string

const (
	ProviderReferenceType = ReferenceType("Provider")
	FischerReferenceType  = ReferenceType("Fischer")
)

type ProviderSpec struct {
	// A name of another provider or fischer, depending on the reference type.
	Reference string `json:"reference,omitempty" protobuf:"bytes,1,opt,name=reference"`
	// The reference type, defaults to "Provider" if reference is set.
	ReferenceType *ReferenceType `json:"referenceType,omitempty" protobuf:"bytes,2,opt,name=referenceType"`

	DownloadName string  `json:"downloadName" protobuf:"bytes,3,name=downloadName"`
	Version      *string `json:"version,omitempty" protobuf:"bytes,4,opt,name=version"`
	DownloadUrl  *string `json:"downloadUrl,omitempty" protobuf:"bytes,5,opt,name=downloadUrl"`
}

type ProviderStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:prerelease-lifecycle-gen:introduced=1.0
// +k8s:prerelease-lifecycle-gen:removed=1.10

type Provider struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   ProviderSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status string       `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}
