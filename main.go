/*
Copyright 2016 The Kubernetes Authors.

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

package main

import (
	"os"

	"hykube.io/apiserver/pkg/cmd/server"
	genericapiserver "k8s.io/apiserver/pkg/server"
	"k8s.io/component-base/cli"
)

// https://github.com/kubernetes/sample-apiserver?tab=readme-ov-file#authentication-plugins
import _ "k8s.io/client-go/plugin/pkg/client/auth"

func main() {
	ctx := genericapiserver.SetupSignalContext()
	options := server.NewHykubeServerOptions(os.Stdout, os.Stderr)
	cmd := server.NewCommandStartHykubeServer(ctx, options)
	code := cli.Run(cmd)
	os.Exit(code)
}
