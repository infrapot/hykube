/*
Copyright 2024 by infrapot

This program is a free software product. You can redistribute it and/or
modify it under the terms of the GNU Affero General Public License (AGPL)
version 3 as published by the Free Software Foundation.

For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
*/

package v1alpha1

import (
	"github.com/infrapot/hykube/pkg/apis/hykube"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:prerelease-lifecycle-gen:introduced=1.0
// +k8s:prerelease-lifecycle-gen:removed=1.10

// ProviderList is a list of Provider objects.
type ProviderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []Provider `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type ProviderSpec struct {
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

	Spec     ProviderSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status   string       `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
	Filename string       `json:"filename,omitempty" protobuf:"bytes,4,opt,name=filename"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:prerelease-lifecycle-gen:introduced=1.0
// +k8s:prerelease-lifecycle-gen:removed=1.10

// PlanList is a list of Plan objects.
type PlanList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Items []Plan `json:"items" protobuf:"bytes,2,rep,name=items"`
}

type ProviderPlanDetails struct {
	Name   string
	Config interface{}
}

func (p *ProviderPlanDetails) DeepCopy() *ProviderPlanDetails {
	ret := ProviderPlanDetails{}

	ret.Name = p.Name
	ret.Config = p.Config // TODO: we assume configuration is always immutable; maybe marshal and unmarshal is better

	return &ret
}

type PlanSpec struct {
	Provider hykube.ProviderPlanDetails
}

type PlanStatus struct {
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +k8s:prerelease-lifecycle-gen:introduced=1.0
// +k8s:prerelease-lifecycle-gen:removed=1.10

type Plan struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	Spec   PlanSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
	Status string   `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}
