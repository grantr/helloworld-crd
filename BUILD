load("@io_bazel_rules_go//go:def.bzl", "gazelle", "go_binary", "go_library", "go_prefix")

go_prefix("k8s.io/sample-controller")

gazelle(
    name = "gazelle",
    external = "vendored",
)

go_library(
    name = "go_default_library",
    srcs = ["main.go"],
    importpath = "k8s.io/sample-controller",
    visibility = ["//visibility:private"],
    deps = [
        "//pkg/client/clientset/versioned:go_default_library",
        "//pkg/client/informers/externalversions:go_default_library",
        "//pkg/controller:go_default_library",
        "//pkg/controller/bar:go_default_library",
        "//pkg/controller/foo:go_default_library",
        "//pkg/signals:go_default_library",
        "//vendor/github.com/golang/glog:go_default_library",
        "//vendor/k8s.io/client-go/informers:go_default_library",
        "//vendor/k8s.io/client-go/kubernetes:go_default_library",
        "//vendor/k8s.io/client-go/tools/clientcmd:go_default_library",
    ],
)

go_binary(
    name = "sample-controller",
    embed = [":go_default_library"],
    importpath = "k8s.io/sample-controller",
    pure = "on",
    visibility = ["//visibility:public"],
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
        "//artifacts/examples:authz",
        "//artifacts/examples:foo",
        "//artifacts/examples:bar",
        ":controller",
    ],
)
