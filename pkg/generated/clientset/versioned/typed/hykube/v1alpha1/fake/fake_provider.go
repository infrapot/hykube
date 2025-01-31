/*
Copyright 2024 by infrapot

This program is a free software product. You can redistribute it and/or
modify it under the terms of the GNU Affero General Public License (AGPL)
version 3 as published by the Free Software Foundation.

For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
*/
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"
	json "encoding/json"
	"fmt"

	v1alpha1 "github.com/infrapot/hykube/pkg/apis/hykube/v1alpha1"
	hykubev1alpha1 "github.com/infrapot/hykube/pkg/generated/applyconfiguration/hykube/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeProviders implements ProviderInterface
type FakeProviders struct {
	Fake *FakeHykubeV1alpha1
	ns   string
}

var providersResource = v1alpha1.SchemeGroupVersion.WithResource("providers")

var providersKind = v1alpha1.SchemeGroupVersion.WithKind("Provider")

// Get takes name of the provider, and returns the corresponding provider object, and an error if there is any.
func (c *FakeProviders) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.Provider, err error) {
	emptyResult := &v1alpha1.Provider{}
	obj, err := c.Fake.
		Invokes(testing.NewGetActionWithOptions(providersResource, c.ns, name, options), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Provider), err
}

// List takes label and field selectors, and returns the list of Providers that match those selectors.
func (c *FakeProviders) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.ProviderList, err error) {
	emptyResult := &v1alpha1.ProviderList{}
	obj, err := c.Fake.
		Invokes(testing.NewListActionWithOptions(providersResource, providersKind, c.ns, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.ProviderList{ListMeta: obj.(*v1alpha1.ProviderList).ListMeta}
	for _, item := range obj.(*v1alpha1.ProviderList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested providers.
func (c *FakeProviders) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchActionWithOptions(providersResource, c.ns, opts))

}

// Create takes the representation of a provider and creates it.  Returns the server's representation of the provider, and an error, if there is any.
func (c *FakeProviders) Create(ctx context.Context, provider *v1alpha1.Provider, opts v1.CreateOptions) (result *v1alpha1.Provider, err error) {
	emptyResult := &v1alpha1.Provider{}
	obj, err := c.Fake.
		Invokes(testing.NewCreateActionWithOptions(providersResource, c.ns, provider, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Provider), err
}

// Update takes the representation of a provider and updates it. Returns the server's representation of the provider, and an error, if there is any.
func (c *FakeProviders) Update(ctx context.Context, provider *v1alpha1.Provider, opts v1.UpdateOptions) (result *v1alpha1.Provider, err error) {
	emptyResult := &v1alpha1.Provider{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateActionWithOptions(providersResource, c.ns, provider, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Provider), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeProviders) UpdateStatus(ctx context.Context, provider *v1alpha1.Provider, opts v1.UpdateOptions) (result *v1alpha1.Provider, err error) {
	emptyResult := &v1alpha1.Provider{}
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceActionWithOptions(providersResource, "status", c.ns, provider, opts), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Provider), err
}

// Delete takes name of the provider and deletes it. Returns an error if one occurs.
func (c *FakeProviders) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteActionWithOptions(providersResource, c.ns, name, opts), &v1alpha1.Provider{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeProviders) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionActionWithOptions(providersResource, c.ns, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.ProviderList{})
	return err
}

// Patch applies the patch and returns the patched provider.
func (c *FakeProviders) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.Provider, err error) {
	emptyResult := &v1alpha1.Provider{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(providersResource, c.ns, name, pt, data, opts, subresources...), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Provider), err
}

// Apply takes the given apply declarative configuration, applies it and returns the applied provider.
func (c *FakeProviders) Apply(ctx context.Context, provider *hykubev1alpha1.ProviderApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Provider, err error) {
	if provider == nil {
		return nil, fmt.Errorf("provider provided to Apply must not be nil")
	}
	data, err := json.Marshal(provider)
	if err != nil {
		return nil, err
	}
	name := provider.Name
	if name == nil {
		return nil, fmt.Errorf("provider.Name must be provided to Apply")
	}
	emptyResult := &v1alpha1.Provider{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(providersResource, c.ns, *name, types.ApplyPatchType, data, opts.ToPatchOptions()), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Provider), err
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *FakeProviders) ApplyStatus(ctx context.Context, provider *hykubev1alpha1.ProviderApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.Provider, err error) {
	if provider == nil {
		return nil, fmt.Errorf("provider provided to Apply must not be nil")
	}
	data, err := json.Marshal(provider)
	if err != nil {
		return nil, err
	}
	name := provider.Name
	if name == nil {
		return nil, fmt.Errorf("provider.Name must be provided to Apply")
	}
	emptyResult := &v1alpha1.Provider{}
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceActionWithOptions(providersResource, c.ns, *name, types.ApplyPatchType, data, opts.ToPatchOptions(), "status"), emptyResult)

	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v1alpha1.Provider), err
}
