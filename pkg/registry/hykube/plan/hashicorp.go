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
	"github.com/go-resty/resty/v2"
	"github.com/infrapot/hykube/pkg/apis/hykube/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"time"
)

type hashicorpPlan struct {
	client *resty.Client
	store  *genericregistry.Store
}

const dataDir = "/data"

func (h *hashicorpPlan) Plan(ctx context.Context, plan *v1alpha1.Plan) error {
	//config, err := rest.InClusterConfig()
	//if err != nil {
	//	return err
	//}
	//
	//res, err := h.getResources(ctx, config)
	//if err != nil {
	//	return err
	//}

	//res.Object

	return nil
}

func (h *hashicorpPlan) getResources(ctx context.Context, config *rest.Config) (*unstructured.UnstructuredList, error) {
	client, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("cannot get k8s config : %w", err)
	}

	res, err := client.Resource(schema.GroupVersionResource{
		Group:    "aws.hykube.io",
		Version:  "v1",
		Resource: "aws-s3-buckets",
	}).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("cannot get resources : %w", err)
	}
	return res, nil
}

type HashicorpPlan interface {
	Plan(ctx context.Context, provider *v1alpha1.Plan) error
}

func NewHashicorpPlan(store *genericregistry.Store) HashicorpPlan {
	client := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second)
	return &hashicorpPlan{
		client: client,
		store:  store,
	}
}
