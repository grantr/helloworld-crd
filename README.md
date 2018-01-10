# sample-controller

This repository implements a simple controller for watching Foo and Bar resources as
defined with a CustomResourceDefinition (CRD).

This particular example demonstrates how to perform basic operations such as:

* How to register two new custom resource (custom resource type) of type `Foo` and `Bar` using a CustomResourceDefinition.
* How to create/get/list instances of your new resource types.
* How to setup a controller on resource handling create/update/delete events.

## Purpose

This is an example of how to build a kube-like controller with a single type.

## Updating Types 

This makes use of the generators in [k8s.io/code-generator](https://github.com/kubernetes/code-generator)
to generate a typed client, informers, listers and deep-copy functions. You can
do this yourself using the `./hack/update-codegen.sh` script.

The `update-codegen` script will automatically generate the following files &
directories:

* `pkg/apis/samplecontroller/v1alpha1/zz_generated.deepcopy.go`
* `pkg/client/`

Changes should not be made to these files manually, and when creating your own
controller based off of this implementation you should not copy these files and
instead run the `update-codegen` script to generate your own.


## Updating Deps

As new external dependencies are added, they will need to be vendored using `dep`.
To manage dependency updates, we have a script `./hack/update-deps.sh` which will
run `dep`, and then run `Gazelle` to generate `BUILD` files for them.


## Running Locally

### One-time setup

We need to register our custom resource once:
```shell
kubectl create -f artifacts/examples/foo.yaml
kubectl create -f artifacts/examples/bar.yaml
```

### Development

To run the controller locally against the current K8s cluster context, run:

```sh
# assumes you have a working kubeconfig, not required if operating in-cluster
go run *.go -kubeconfig=$HOME/.kube/config -logtostderr=true -stderrthreshold=INFO

# create a custom resource of type Foo
kubectl create -f artifacts/examples/example-foo.yaml

# create a custom resource of type Bar
kubectl create -f artifacts/examples/example-bar.yaml

```

### Cleanup

You can clean up the created CustomResourceDefinition with:

```shell
kubectl delete crd foos.samplecontroller.k8s.io
kubectl delete crd bars.samplecontroller.k8s.io
```

## Running On-Cluster

### One-time setup

To tell Bazel where to publish images, and to which cluster to deploy:

```shell
# You can put these definitions in .bashrc, so this is one-time setup.
export DOCKER_REPO_OVERRIDE=us.gcr.io/project
# See: kubectl config get-contexts
export K8S_CLUSTER_OVERRIDE=cluster-name

# Forces Bazel to pick up these changes (don't put in .bashrc)
bazel clean
```

Note that this expects your Docker authorization is [properly configured](
https://github.com/bazelbuild/rules_docker#authorization).

### Standing it up

You can stand up a version of this controller on-cluster with:
```shell
# This will register the CRD and deploy the controller to start acting on them.
bazel run :everything.create
```

To test things out, you can create an example `Foo` with:
```shell
bazel run artifacts/examples:example-foo.create
# Or a Bar!
bazel run artifacts/examples:example-bar.create
```

### Iterating

As you make changes to the code, you can redeploy your controller with:
```shell
bazel run :controller.replace
```

Two things of note:
1. If your (external) dependencies have changed, you should: `./hack/update-deps.sh`.
1. If your type definitions have changed, you should: `./hack/update-codegen.sh`.

If only internal dependencies have changed, and you want to avoid the `dep`
portion of `./hack/update-deps.sh`, you can just run `Gazelle` with:
```shell
bazel run //:gazelle -- -proto=disable
```

### Cleanup

You can clean up everything with:
```shell
bazel run :everything.delete
```

## Use Cases

CustomResourceDefinitions can be used to implement custom resource types for your Kubernetes cluster.
These act like most other Resources in Kubernetes, and may be `kubectl apply`'d, etc.

Some example use cases:

* Provisioning/Management of external datastores/databases (eg. CloudSQL/RDS instances)
* Higher level abstractions around Kubernetes primitives (eg. a single Resource to define an etcd cluster, backed by a Service and a ReplicationController)

## Defining types

Each instance of your custom resource has an attached Spec, which should be defined via a `struct{}` to provide data format validation.
In practice, this Spec is arbitrary key-value data that specifies the configuration/behavior of your Resource.

For example, if you were implementing a custom resource for a Database, you might provide a DatabaseSpec like the following:

``` go
type DatabaseSpec struct {
	Databases []string `json:"databases"`
	Users     []User   `json:"users"`
	Version   string   `json:"version"`
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
```

## Compatibility

HEAD of this repository will match HEAD of k8s.io/apimachinery and
k8s.io/client-go.

## Where does it come from?

`sample-controller` is synced from
https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/sample-controller.
Code changes are made in that location, merged into k8s.io/kubernetes and
later synced here.
