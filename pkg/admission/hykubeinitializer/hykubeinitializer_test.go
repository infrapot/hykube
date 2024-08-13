/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package hykubeinitializer_test

import (
	"context"
	"testing"
	"time"

	"hykube.io/apiserver/pkg/admission/hykubeinitializer"
	"hykube.io/apiserver/pkg/generated/clientset/versioned/fake"
	informers "hykube.io/apiserver/pkg/generated/informers/externalversions"
	"k8s.io/apiserver/pkg/admission"
)

// TestWantsInternalHykubeInformerFactory ensures that the informer factory is injected
// when the wantInternalHykubeInformerFactory interface is implemented by a plugin.
func TestWantsInternalHykubeInformerFactory(t *testing.T) {
	cs := &fake.Clientset{}
	sf := informers.NewSharedInformerFactory(cs, time.Duration(1)*time.Second)
	target := hykubeinitializer.New(sf)

	wantHykubeInformerFactory := &wantInternalHykubeInformerFactory{}
	target.Initialize(wantHykubeInformerFactory)
	if wantHykubeInformerFactory.sf != sf {
		t.Errorf("expected informer factory to be initialized")
	}
}

// wantInternalHykubeInformerFactory is a test stub that fulfills the WantsInternalHykubeInformerFactory interface
type wantInternalHykubeInformerFactory struct {
	sf informers.SharedInformerFactory
}

func (f *wantInternalHykubeInformerFactory) SetInternalHykubeInformerFactory(sf informers.SharedInformerFactory) {
	f.sf = sf
}
func (f *wantInternalHykubeInformerFactory) Admit(ctx context.Context, a admission.Attributes, o admission.ObjectInterfaces) error {
	return nil
}
func (f *wantInternalHykubeInformerFactory) Handles(o admission.Operation) bool { return false }
func (f *wantInternalHykubeInformerFactory) ValidateInitialization() error      { return nil }

var _ admission.Interface = &wantInternalHykubeInformerFactory{}
var _ hykubeinitializer.WantsInternalHykubeInformerFactory = &wantInternalHykubeInformerFactory{}