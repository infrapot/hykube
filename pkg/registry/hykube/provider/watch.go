/*
 * Copyright 2024 by infrapot
 *
 * This program is a free software product. You can redistribute it and/or
 * modify it under the terms of the GNU Affero General Public License (AGPL)
 * version 3 as published by the Free Software Foundation.
 *
 * For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
 */

package provider

import (
	"context"
	"hykube.io/apiserver/pkg/apis/hykube/v1alpha1"
	watch2 "k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/klog/v2"
)

type watcher struct {
	watch             watch2.Interface
	hashicorpProvider HashicorpProvider
}

func (p *watcher) Start() {

	go func() {
		for {
			//var event *watch2.Event
			event := <-p.watch.ResultChan()
			if event.Type == watch2.Added {
				provider := event.Object.(*v1alpha1.Provider)

				ctx := request.WithNamespace(context.TODO(), provider.Namespace)

				klog.Infof("Found added event for: %v", provider.Name)

				err := p.hashicorpProvider.Initialize(ctx, provider)

				if err != nil {
					klog.ErrorS(err, "Couldn't initialize provider: %s", provider.Name)
				}
			}
		}
	}()
}
