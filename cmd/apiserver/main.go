/*
Copyright 2024 by infrapot

This program is a free software product. You can redistribute it and/or
modify it under the terms of the GNU Affero General Public License (AGPL)
version 3 as published by the Free Software Foundation.

For details, see the GNU AGPL at: http://www.gnu.org/licenses/agpl-3.0.html
*/

package main

import (
	"k8s.io/klog"
	"os"
	"sigs.k8s.io/apiserver-runtime/pkg/builder"
	// +kubebuilder:scaffold:resource-imports
)

func main() {
	os.Setenv("KUBERNETES_SERVICE_HOST", "0.0.0.0")
	os.Setenv("KUBERNETES_SERVICE_PORT", "443")
	err := builder.APIServer.
		DisableAuthorization().
		WithOptionsFns(func(options *builder.ServerOptions) *builder.ServerOptions {
			options.RecommendedOptions.CoreAPI = nil
			options.RecommendedOptions.Admission = nil
			return options
		}).
		WithLocalDebugExtension().
		WithoutEtcd().
		Execute()
	if err != nil {
		klog.Fatal(err)
	}
}
