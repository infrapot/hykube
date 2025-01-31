/*
Copyright 2024 by infrapot

This program is a free software product. You can redistribute it and/or
modify it under the terms of the GNU Affero General Public License (AGPL)
version 3 as published by the Free Software Foundation.

For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
*/
// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1alpha1

// ProviderSpecApplyConfiguration represents a declarative configuration of the ProviderSpec type for use
// with apply.
type ProviderSpecApplyConfiguration struct {
	DownloadName *string `json:"downloadName,omitempty"`
	Version      *string `json:"version,omitempty"`
	DownloadUrl  *string `json:"downloadUrl,omitempty"`
}

// ProviderSpecApplyConfiguration constructs a declarative configuration of the ProviderSpec type for use with
// apply.
func ProviderSpec() *ProviderSpecApplyConfiguration {
	return &ProviderSpecApplyConfiguration{}
}

// WithDownloadName sets the DownloadName field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DownloadName field is set to the value of the last call.
func (b *ProviderSpecApplyConfiguration) WithDownloadName(value string) *ProviderSpecApplyConfiguration {
	b.DownloadName = &value
	return b
}

// WithVersion sets the Version field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Version field is set to the value of the last call.
func (b *ProviderSpecApplyConfiguration) WithVersion(value string) *ProviderSpecApplyConfiguration {
	b.Version = &value
	return b
}

// WithDownloadUrl sets the DownloadUrl field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the DownloadUrl field is set to the value of the last call.
func (b *ProviderSpecApplyConfiguration) WithDownloadUrl(value string) *ProviderSpecApplyConfiguration {
	b.DownloadUrl = &value
	return b
}
