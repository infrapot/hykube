//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2024 by infrapot

This program is a free software product. You can redistribute it and/or
modify it under the terms of the GNU Affero General Public License (AGPL)
version 3 as published by the Free Software Foundation.

For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
*/
// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha1

import (
	hykube "hykube.io/apiserver/pkg/apis/hykube"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*Provider)(nil), (*hykube.Provider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_Provider_To_hykube_Provider(a.(*Provider), b.(*hykube.Provider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*hykube.Provider)(nil), (*Provider)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_hykube_Provider_To_v1alpha1_Provider(a.(*hykube.Provider), b.(*Provider), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ProviderList)(nil), (*hykube.ProviderList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ProviderList_To_hykube_ProviderList(a.(*ProviderList), b.(*hykube.ProviderList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*hykube.ProviderList)(nil), (*ProviderList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_hykube_ProviderList_To_v1alpha1_ProviderList(a.(*hykube.ProviderList), b.(*ProviderList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*hykube.ProviderSpec)(nil), (*ProviderSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_hykube_ProviderSpec_To_v1alpha1_ProviderSpec(a.(*hykube.ProviderSpec), b.(*ProviderSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddConversionFunc((*ProviderSpec)(nil), (*hykube.ProviderSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha1_ProviderSpec_To_hykube_ProviderSpec(a.(*ProviderSpec), b.(*hykube.ProviderSpec), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha1_Provider_To_hykube_Provider(in *Provider, out *hykube.Provider, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1alpha1_ProviderSpec_To_hykube_ProviderSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	out.Status = in.Status
	return nil
}

// Convert_v1alpha1_Provider_To_hykube_Provider is an autogenerated conversion function.
func Convert_v1alpha1_Provider_To_hykube_Provider(in *Provider, out *hykube.Provider, s conversion.Scope) error {
	return autoConvert_v1alpha1_Provider_To_hykube_Provider(in, out, s)
}

func autoConvert_hykube_Provider_To_v1alpha1_Provider(in *hykube.Provider, out *Provider, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_hykube_ProviderSpec_To_v1alpha1_ProviderSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	out.Status = in.Status
	return nil
}

// Convert_hykube_Provider_To_v1alpha1_Provider is an autogenerated conversion function.
func Convert_hykube_Provider_To_v1alpha1_Provider(in *hykube.Provider, out *Provider, s conversion.Scope) error {
	return autoConvert_hykube_Provider_To_v1alpha1_Provider(in, out, s)
}

func autoConvert_v1alpha1_ProviderList_To_hykube_ProviderList(in *ProviderList, out *hykube.ProviderList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]hykube.Provider, len(*in))
		for i := range *in {
			if err := Convert_v1alpha1_Provider_To_hykube_Provider(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1alpha1_ProviderList_To_hykube_ProviderList is an autogenerated conversion function.
func Convert_v1alpha1_ProviderList_To_hykube_ProviderList(in *ProviderList, out *hykube.ProviderList, s conversion.Scope) error {
	return autoConvert_v1alpha1_ProviderList_To_hykube_ProviderList(in, out, s)
}

func autoConvert_hykube_ProviderList_To_v1alpha1_ProviderList(in *hykube.ProviderList, out *ProviderList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Provider, len(*in))
		for i := range *in {
			if err := Convert_hykube_Provider_To_v1alpha1_Provider(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_hykube_ProviderList_To_v1alpha1_ProviderList is an autogenerated conversion function.
func Convert_hykube_ProviderList_To_v1alpha1_ProviderList(in *hykube.ProviderList, out *ProviderList, s conversion.Scope) error {
	return autoConvert_hykube_ProviderList_To_v1alpha1_ProviderList(in, out, s)
}

func autoConvert_v1alpha1_ProviderSpec_To_hykube_ProviderSpec(in *ProviderSpec, out *hykube.ProviderSpec, s conversion.Scope) error {
	// WARNING: in.Reference requires manual conversion: does not exist in peer-type
	// WARNING: in.ReferenceType requires manual conversion: inconvertible types (*hykube.io/apiserver/pkg/apis/hykube/v1alpha1.ReferenceType vs hykube.io/apiserver/pkg/apis/hykube.ReferenceType)
	return nil
}

func autoConvert_hykube_ProviderSpec_To_v1alpha1_ProviderSpec(in *hykube.ProviderSpec, out *ProviderSpec, s conversion.Scope) error {
	// WARNING: in.ProviderReference requires manual conversion: does not exist in peer-type
	// WARNING: in.ReferenceType requires manual conversion: inconvertible types (hykube.io/apiserver/pkg/apis/hykube.ReferenceType vs *hykube.io/apiserver/pkg/apis/hykube/v1alpha1.ReferenceType)
	return nil
}
