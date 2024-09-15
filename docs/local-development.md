# Local development

This document will take you through setting up and trying the sample apiserver on a local k8s from a fresh clone of this repo.

## Pre requisites

- K8S version **1.30** + Docker
- Go **1.22** or later

## Code generation

If you change the API object type definitions in `pkg/apis/.../types.go` files then you will need to update the files generated from the type definitions. 

To do this, first call `go mod vendor` to get correct vendored deps and then invoke `./hack/update-codegen.sh` with `hykube` as your current working directory; the script takes no arguments.

## Build the binary

If you work on macOS, invoke the following command first to install GCC:
```shell
brew tap SergioBenitez/osxct
brew install x86_64-unknown-linux-gnu
```

Next we will want to create a new binary to both test we can build the server and to use for the container image.

From the root of this repo, where ```main.go``` is located, run the following command (on macOS):
```shell
CC=x86_64-unknown-linux-gnu-gcc CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -a -o artifacts/apiserver-image/hykube-apiserver
```

On Linux, run the following command:
```shell
CGO_ENABLED=1 GOARCH=amd64 GOOS=linux go build -a -o artifacts/apiserver-image/hykube-apiserver
```

if everything went well, you should have a binary called `hykube-apiserver` present in `artifacts/simple-image`.

## Build the container image

Using the binary we just built, we will now create a Docker image and push it to our Dockerhub registry so that we deploy it to our cluster.
There is a sample `Dockerfile` located in `artifacts/apiserver-image` we will use this to build our own image.

Again from the root of this repo run the following commands:
```shell
docker build -t hykube-apiserver:latest ./artifacts/apiserver-image
```

## Deploy to K8S

We will need to create several objects in order to setup the API server, so you will need to ensure you have the `kubectl` tool installed. [Install kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl/).

```shell
# create the namespace to run the apiserver in
kubectl create ns hykube

# create the service account used to run the server
kubectl create -f artifacts/deployment/sa.yaml -n hykube

# create the rolebindings that allow the service account user to delegate authz back to the kubernetes master for incoming requests to the apiserver
kubectl create -f artifacts/deployment/auth-delegator.yaml -n kube-system
kubectl create -f artifacts/deployment/auth-reader.yaml -n kube-system

# create rbac roles and clusterrolebinding that allow the service account user to use admission webhooks
kubectl create -f artifacts/deployment/rbac.yaml
kubectl create -f artifacts/deployment/rbac-bind.yaml

# create the, PV, PVC, service and replication controller
kubectl create -f artifacts/deployment/pv.yaml -n hykube
kubectl create -f artifacts/deployment/pvc.yaml -n hykube
kubectl create -f artifacts/deployment/deployment.yaml -n hykube
kubectl create -f artifacts/deployment/service.yaml -n hykube

# create the apiservice object that tells kubernetes about your api extension and where in the cluster the server is located
kubectl create -f artifacts/deployment/apiservice.yaml
```

## Test that your setup has worked

You should now be able to create the resource type `Provider` which is the resource type registered by the API server.

```shell
kubectl create -f artifacts/providers/aws.yaml
# provider "aws" created
```

You can then get this resource by running:

```shell
kubectl get provider aws -o custom-columns=Name:.metadata.name,status:.status,CreatedAt:.metadata.creationTimestamp
# Name   status        CreatedAt
# aws    adding CRDs   2024-09-12T09:15:29Z
```

After the provider is initialized and all ~1500 CRDs are added (can take a while), you should see status `ready`:
```shell
kubectl get provider aws -o custom-columns=Name:.metadata.name,status:.status,CreatedAt:.metadata.creationTimestamp
# Name   status        CreatedAt
# aws    ready         2024-09-12T09:15:29Z
```

## Deploy locally resources

Once the provider is ready, you can add a test AWS S3 bucket:

```shell
kubectl create -f artifacts/aws-test/s3-bucket.yaml                                                                       
kubectl delete -f artifacts/aws-test/s3-bucket.yaml                                                                       
# aws-s3-bucket.aws.hykube.io/test-bucket created
```

## Local binary run

Generate certificates for CA and test user `development` in the superuser group `system:masters` for local testing:
```shell
openssl req -nodes -new -x509 -keyout ca.key -out ca.crt

openssl req -out client.csr -new -newkey rsa:4096 -nodes -keyout client.key -subj "/CN=development/O=system:masters"
openssl x509 -req -days 365 -in client.csr -CA ca.crt -CAkey ca.key -set_serial 01 -sha256 -out client.crt

openssl pkcs12 -export -in ./client.crt -inkey ./client.key -out client.p12 -passout pass:password
```
Set the generated certs and server config:
```shell
kubectl config set --kubeconfig ./config clusters.hykube.certificate-authority-data $(cat ca.crt | base64 -i -)
kubectl config set --kubeconfig ./config users.hykube.client-certificate-data $(cat client.crt | base64 -i -)
kubectl config set --kubeconfig ./config users.hykube.client-key-data $(cat client.key | base64 -i -)
```
Add a context to kubeconfig:
```yaml
contexts:
- context:
    cluster: hykube
    user: hykube
  name: hykube
current-context: hykube
```
Set test server address:
```yaml
clusters:
- cluster:
    certificate-authority-data: [SET ABOVE]
    server: https://127.0.0.1:8443
  name: hykube
```

Run binary with following arguments:
```shell
hykube-apiserver --secure-port 8443 --v=7 \
   --client-ca-file ca.crt \
   --kubeconfig ./config \
   --authentication-kubeconfig ./config \
   --authorization-kubeconfig ./config \
   --authentication-skip-lookup
```

Issue a sample call to API Server:
```shell
wget -O- --no-check-certificate --certificate client.crt --private-key client.key \
https://localhost:8443/apis/hykube.io/v1alpha1/namespaces/default/providers
```