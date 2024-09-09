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
	"archive/zip"
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/terraform/configs/configschema"
	"github.com/hashicorp/terraform/providers"
	"hykube.io/apiserver/pkg/apis/hykube"
	"hykube.io/apiserver/pkg/apis/hykube/v1alpha1"
	"io"
	apiextensionv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	crd_clientset "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	watch2 "k8s.io/apimachinery/pkg/watch"
	"k8s.io/apiserver/pkg/endpoints/request"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	tfplugin "github.com/hashicorp/terraform/plugin"
)

type watcher struct {
	watch watch2.Interface
	store *genericregistry.Store
}

func (p *watcher) Start() {
	client := resty.New().
		SetRetryCount(3).
		SetRetryWaitTime(1 * time.Second)
	go func() {
		for {
			//var event *watch2.Event
			event := <-p.watch.ResultChan()
			if event.Type == watch2.Added {
				provider := event.Object.(*v1alpha1.Provider)

				ctx := request.WithNamespace(context.TODO(), provider.Namespace)

				klog.Infof("Found added event for: %s", provider.Name)

				key, err := p.store.KeyFunc(ctx, provider.Name)
				if err != nil {
					klog.ErrorS(err, "cannot generate key")
					continue
				}

				preconditions := storage.Preconditions{UID: &provider.UID, ResourceVersion: &provider.ResourceVersion}

				out := p.store.NewFunc()

				err = p.updateProviderStatus(ctx, key, out, preconditions, "downloading provider")
				if err != nil {
					klog.ErrorS(err, "cannot update provider")
					continue
				}

				var fullURLFile string
				if provider.Spec.DownloadUrl != nil {
					fullURLFile = *provider.Spec.DownloadUrl
				} else {
					fullURLFile = "https://releases.hashicorp.com/" + provider.Spec.DownloadName + "/" +
						*provider.Spec.Version + "/" + provider.Spec.DownloadName + "_" + *provider.Spec.Version +
						"_linux_amd64.zip" // TODO: detect latest version by default
				}

				filename := provider.Name + ".zip"
				_, err = client.R().
					SetOutput(filename).
					Get(fullURLFile)
				if err != nil {
					klog.ErrorS(err, "Couldn't download file")
					return
				}

				klog.Infof("Downloaded a provider from: %s", fullURLFile)
				err = p.updateProviderStatus(ctx, key, out, preconditions, "downloaded provider")
				if err != nil {
					klog.ErrorS(err, "cannot update provider")
					continue
				}

				err = p.updateProviderStatus(ctx, key, out, preconditions, "extracting provider")
				if err != nil {
					klog.ErrorS(err, "cannot update provider")
					continue
				}

				providerFilename, err := p.extractFile(err, filename)
				if err != nil {
					klog.ErrorS(err, "Couldn't extract file")
					return
				}

				err = p.updateProviderStatus(ctx, key, out, preconditions, "getting provider schema")
				if err != nil {
					klog.ErrorS(err, "cannot update provider")
					continue
				}

				schema, err := p.getProviderSchema(providerFilename, false)
				if err != nil {
					klog.ErrorS(err, "Couldn't get provider schema")
					return
				}

				err = p.updateProviderStatus(ctx, key, out, preconditions, "adding CRDs")
				if err != nil {
					klog.ErrorS(err, "cannot update provider")
					continue
				}

				err = p.addCDRs(schema, provider)
				if err != nil {
					klog.ErrorS(err, "Couldn't get provider schema")
					return
				}

				klog.Infof("Added a provider from: %s", fullURLFile)
			}
		}
	}()
}

func (p *watcher) updateProviderStatus(ctx context.Context, key string, out runtime.Object, preconditions storage.Preconditions, statusValue string) error {
	return p.store.Storage.GuaranteedUpdate(
		ctx, key, out, false, &preconditions,
		storage.SimpleUpdate(func(existing runtime.Object) (runtime.Object, error) {
			existingProvider, ok := existing.(*hykube.Provider)
			if !ok {
				// wrong type
				return nil, fmt.Errorf("expected *hykube.Provider, got %v", existing)
			}
			existingProvider.Status = statusValue
			return existingProvider, nil
		}),
		false, // watcher doesn't get notified if it's dry run
		nil,
	)
}

func (p *watcher) extractFile(err error, filename string) (string, error) {
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

func (p *watcher) getProviderSchema(providerFilename string, verbose bool) (*providers.GetSchemaResponse, error) {
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

func (p *watcher) addCDRs(schemaResponse *providers.GetSchemaResponse, provider *v1alpha1.Provider) error {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err
	}

	crdClient, err := crd_clientset.NewForConfig(config)
	CRDGroup := provider.Name + ".hykube.io"
	CRDVersion := "v1alpha1"

	if err != nil {
		return err
	}
	for k, v := range schemaResponse.ResourceTypes {
		// TODO look at ImpliedType and blockTypes
		kind := strings.ToLower(k)
		CRDPlural := kind + "s" // TODO check if it's a standard way of doing it
		FullCRDName := CRDPlural + "." + CRDGroup

		var requiredFields []string
		properties := make(map[string]apiextensionv1.JSONSchemaProps, len(v.Block.Attributes))

		for attrK, attrV := range v.Block.Attributes {
			if attrV.Required {
				requiredFields = append(requiredFields, attrK)
			}
			properties[attrK] = apiextensionv1.JSONSchemaProps{
				Type: p.attributeType(attrV),
			}
		}

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
								Type:       "object",
								Required:   requiredFields,
								Properties: properties,
							},
						},
					},
				},
				Scope: apiextensionv1.NamespaceScoped,
				Names: apiextensionv1.CustomResourceDefinitionNames{
					Plural: CRDPlural,
					Kind:   kind,
				},
			},
		}

		_, err = crdClient.ApiextensionsV1().CustomResourceDefinitions().Create(context.TODO(), crd, meta_v1.CreateOptions{})
		if err != nil && apierrors.IsAlreadyExists(err) {
			return nil
		}
	}

	return nil
}

func (p *watcher) attributeType(attrV *configschema.Attribute) string {
	if attrV.Type.IsPrimitiveType() {
		return attrV.Type.FriendlyName()
	} else { // TODO improve complex type
		if attrV.Type.IsListType() {
			return "list"
		} else {
			return "object"
		}
	}
}
