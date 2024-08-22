/*
 * Copyright 2024 by infrapot
 *
 * This program is a free software product. You can redistribute it and/or
 * modify it under the terms of the GNU Affero General Public License (AGPL)
 * version 3 as published by the Free Software Foundation.
 *
 * For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
 */

package server

import (
	"context"
	"github.com/k3s-io/kine/pkg/endpoint"
	"hykube.io/apiserver/pkg/apiserver"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/runtime/serializer/json"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage/storagebackend"
	"k8s.io/apiserver/pkg/util/flowcontrol/request"
	"time"
)

type proxiedRESTOptionsGetter struct {
	dsn            string
	groupVersioner runtime.GroupVersioner
}

// GetRESTOptions implements RESTOptionsGetter interface.
func (g *proxiedRESTOptionsGetter) GetRESTOptions(resource schema.GroupResource, example runtime.Object) (generic.RESTOptions, error) {

	// A note on etcd connection config - we assume Hykube does not need separate certs for the kine sock connection. If needed, socks connection can be added later on.
	etcdConfig, err := endpoint.Listen(context.TODO(), endpoint.Config{
		Endpoint:       g.dsn,
		NotifyInterval: 1 * time.Second,
	})
	if err != nil {
		return generic.RESTOptions{}, err
	}
	s := json.NewSerializer(json.DefaultMetaFactory, apiserver.Scheme, apiserver.Scheme, false)
	codec := serializer.NewCodecFactory(apiserver.Scheme).
		CodecForVersions(s, s, g.groupVersioner, g.groupVersioner)

	restOptions := generic.RESTOptions{
		ResourcePrefix: resource.String(),
		//Decorator:                 genericregistry.StorageWithCacher(),
		Decorator:                 generic.UndecoratedStorage,
		EnableGarbageCollection:   true,
		DeleteCollectionWorkers:   1,
		CountMetricPollPeriod:     time.Minute,
		StorageObjectCountTracker: request.NewStorageObjectCountTracker(),
		StorageConfig: &storagebackend.ConfigForResource{
			GroupResource: resource,
			Config: storagebackend.Config{
				Prefix: "/kine/",
				Codec:  codec,
				Transport: storagebackend.TransportConfig{
					ServerList:    etcdConfig.Endpoints,
					TrustedCAFile: etcdConfig.TLSConfig.CAFile,
					CertFile:      etcdConfig.TLSConfig.CertFile,
					KeyFile:       etcdConfig.TLSConfig.KeyFile,
				},
			},
		},
	}
	return restOptions, nil
}
