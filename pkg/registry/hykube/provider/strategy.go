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

package provider

import (
	"archive/zip"
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform/providers"
	"hykube.io/apiserver/pkg/apis/hykube"
	"io"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"os"
	"os/exec"
	"strings"
	"time"

	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	tfplugin "github.com/hashicorp/terraform/plugin"
	"hykube.io/apiserver/pkg/apis/hykube/validation"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
)

// NewStrategy creates and returns a providerStrategy instance
func NewStrategy(typer runtime.ObjectTyper) providerStrategy {
	client := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second)

	return providerStrategy{
		typer,
		names.SimpleNameGenerator,
		client,
	}
}

// GetAttrs returns labels.Set, fields.Set, and error in case the given runtime.Object is not a Provider
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, error) {
	apiServer, ok := obj.(*hykube.Provider)
	if !ok {
		return nil, nil, fmt.Errorf("given object is not a Provider")
	}
	return labels.Set(apiServer.ObjectMeta.Labels), SelectableFields(apiServer), nil
}

// MatchProvider is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiServer only interested in specific labels/fields.
func MatchProvider(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *hykube.Provider) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

type providerStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator

	client *resty.Client
}

func (providerStrategy) NamespaceScoped() bool {
	return true
}

func (p providerStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	go func() { // TODO: move to events post created
		//fullURLFile := "https://releases.hashicorp.com/terraform-provider-aws/5.62.0/terraform-provider-aws_5.62.0_linux_amd64.zip" // TODO: detect latest version by default
		//
		//jj, _ := json.Marshal(obj)
		//klog.Info(string(jj))
		//
		//filename := "aws-provider.zip"
		//_, err := p.client.R().
		//	SetOutput(filename).
		//	Get(fullURLFile)
		//if err != nil {
		//	klog.ErrorS(err, "Couldn't download file")
		//	return
		//}
		//
		//klog.Infof("Downloaded a provider from: %s", fullURLFile)
		//
		//providerFilename, err := p.extractFile(err, filename)
		//if err != nil {
		//	klog.ErrorS(err, "Couldn't extract file")
		//	return
		//}
		//
		//schema, err := p.getProviderSchema(providerFilename, false)
		//if err != nil {
		//	klog.ErrorS(err, "Couldn't get provider schema")
		//	return
		//}

		err := p.addCDRs()
		if err != nil {
			klog.ErrorS(err, "Couldn't add CRDs")
			return
		}

	}()
}

func (p providerStrategy) extractFile(err error, filename string) (string, error) {
	r, err := zip.OpenReader(filename)
	if err != nil {
		return "", fmt.Errorf("couldn't open zip file: %w", err)
	}
	defer r.Close()
	providerFilename := ""
	for _, f := range r.File {
		if f.Name == "LICENSE.txt" {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			return "", fmt.Errorf("couldn't open file: %w", err)
		}
		defer rc.Close()
		f, err := os.OpenFile(
			f.Name, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return "", fmt.Errorf("couldn't open new file file: %w", err)
		}
		defer f.Close()
		_, err = io.Copy(f, rc)
		if err != nil {
			return "", fmt.Errorf("couldn't copy file: %w", err)
		}
		providerFilename = f.Name()
	}
	if providerFilename == "" {
		return "", fmt.Errorf("no provider file found")
	}
	return providerFilename, nil
}

func (p providerStrategy) getProviderSchema(providerFilename string, verbose bool) (*providers.GetSchemaResponse, error) {
	options := hclog.LoggerOptions{
		Name:   "plugin",
		Level:  hclog.Error,
		Output: os.Stdout,
	}
	if verbose {
		options.Level = hclog.Trace
	}
	logger := hclog.New(&options)
	client := plugin.NewClient(
		&plugin.ClientConfig{
			Cmd:              exec.Command("/" + providerFilename),
			HandshakeConfig:  tfplugin.Handshake,
			VersionedPlugins: tfplugin.VersionedPlugins,
			Managed:          true,
			Logger:           logger,
			AllowedProtocols: []plugin.Protocol{plugin.ProtocolGRPC},
			AutoMTLS:         true,
		})
	rpcClient, err := client.Client()
	if err != nil {
		return nil, err
	}
	raw, err := rpcClient.Dispense(tfplugin.ProviderPluginName)
	if err != nil {
		return nil, err
	}

	provider := raw.(*tfplugin.GRPCProvider)

	schema := provider.GetSchema()

	client.Kill()

	return &schema, nil
}

func (p providerStrategy) addCDRs() error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	crdClient, err := clientset.NewForConfig(config)

	CRDPlural := strings.ToLower("VPCs")
	CRDGroup := "aws.hykube.io"
	CRDVersion := "v1alpha1"
	FullCRDName := CRDPlural + "." + CRDGroup

	minLen := 1.0
	crd := &apiextensionv1.CustomResourceDefinition{
		ObjectMeta: meta_v1.ObjectMeta{Name: FullCRDName},
		Spec: apiextensionv1.CustomResourceDefinitionSpec{
			Group: CRDGroup,
			Versions: []apiextensionv1.CustomResourceDefinitionVersion{
				{
					Name:    CRDVersion,
					Storage: true,
					Schema: &apiextensionv1.CustomResourceValidation{
						OpenAPIV3Schema: &apiextensionv1.JSONSchemaProps{
							Type:     "object",
							Required: []string{"spec"},
							Properties: map[string]apiextensionv1.JSONSchemaProps{
								"spec": {
									Type:     "object",
									Required: []string{"name"},
									Properties: map[string]apiextensionv1.JSONSchemaProps{
										"name": {
											Type:    "string",
											Minimum: &minLen,
										},
									},
								},
							},
						},
					},
				},
			},
			Scope: apiextensionv1.NamespaceScoped,
			Names: apiextensionv1.CustomResourceDefinitionNames{
				Plural: CRDPlural,
				Kind:   strings.ToLower("VPC"),
			},
		},
	}

	_, err = crdClient.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(), crd, meta_v1.CreateOptions{})
	if err != nil && apierrors.IsAlreadyExists(err) {
		return nil
	}
	return err
}

func (providerStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (providerStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	provider := obj.(*hykube.Provider)
	return validation.ValidateProvider(provider)
}

// WarningsOnCreate returns warnings for the creation of the given object.
func (providerStrategy) WarningsOnCreate(ctx context.Context, obj runtime.Object) []string {
	return nil
}

func (providerStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (providerStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (providerStrategy) Canonicalize(obj runtime.Object) {
}

func (providerStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

// WarningsOnUpdate returns warnings for the given update.
func (providerStrategy) WarningsOnUpdate(ctx context.Context, obj, old runtime.Object) []string {
	return nil
}
