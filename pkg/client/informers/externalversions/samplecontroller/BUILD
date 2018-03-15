load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = ["interface.go"],
    importpath = "k8s.io/sample-controller/pkg/client/informers/externalversions/samplecontroller",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/client/informers/externalversions/internalinterfaces:go_default_library",
        "//pkg/client/informers/externalversions/samplecontroller/v1alpha1:go_default_library",
    ],
)
