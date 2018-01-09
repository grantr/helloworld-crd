load("@io_bazel_rules_go//go:def.bzl", "gazelle", "go_binary", "go_library", "go_prefix")

go_prefix("k8s.io/sample-controller")

gazelle(
    name = "gazelle",
    external = "vendored",
)

go_library(
    name = "go_default_library",
    srcs = [
        "controller.go",
        "main.go",
    ],
    importpath = "k8s.io/sample-controller",
    visibility = ["//visibility:private"],
    deps = [
        "//pkg/client/clientset/versioned:go_default_library",
        "//pkg/client/clientset/versioned/scheme:go_default_library",
        "//pkg/client/informers/externalversions:go_default_library",
        "//pkg/client/listers/samplecontroller/v1alpha1:go_default_library",
        "//pkg/signals:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/api/core/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/api/errors:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/util/wait:go_default_library",
        "//vendor/k8s.io/client-go/informers:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes/scheme:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes/typed/core/v1:go_default_library",
        "//vendor/k8s.io/client-go/tools/cache:go_default_library",
        "//vendor/k8s.io/client-go/tools/clientcmd:go_default_library",
        "//vendor/k8s.io/client-go/tools/record:go_default_library",
        "//vendor/k8s.io/client-go/util/workqueue:go_default_library",
    ],
)

go_binary(
    name = "sample-controller",
    embed = [":go_default_library"],
    importpath = "k8s.io/sample-controller",
    visibility = ["//visibility:public"],
    pure = "on",
)

load("@io_bazel_rules_docker//go:image.bzl", "go_image")

go_image(
    name = "image",
    binary = ":sample-controller",
)

load("@k8s_object//:defaults.bzl", "k8s_object")

k8s_object(
    name = "controller",
    images = {
        "sample-controller:latest": ":image",
    },
    template = "controller.yaml",
)

load("@io_bazel_rules_k8s//k8s:objects.bzl", "k8s_objects")

k8s_objects(
    name = "everything",
    objects = [
        "//artifacts/examples:crd",
        ":controller",
    ],
)
