package(default_visibility = ["//visibility:public"])

load("@k8s_object//:defaults.bzl", "k8s_object")

# TODO(mattmoor): This should move out of examples, it doesn't belong here.
k8s_object(
    name = "crd",
    template = "crd.yaml",
)

k8s_object(
    name = "example",
    template = "example-foo.yaml",
)

k8s_object(
    name = "authz",
    template = "authz.yaml",
)
