load("@io_bazel_rules_go//go:def.bzl", "go_library")

go_library(
    name = "go_default_library",
    srcs = [
        "bar_types.go",
        "doc.go",
        "foo_types.go",
        "register.go",
        "zz_generated.deepcopy.go",
    ],
    importpath = "k8s.io/sample-controller/pkg/apis/samplecontroller/v1alpha1",
    visibility = ["//visibility:public"],
    deps = [
        "//pkg/apis/samplecontroller:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/apis/meta/v1:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime:go_default_library",
        "//vendor/k8s.io/apimachinery/pkg/runtime/schema:go_default_library",
    ],
)
