load("//android:android_grpc_compile.bzl", "android_grpc_compile")
load("@build_bazel_rules_android//android:rules.bzl", "android_library")

def android_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    android_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    android_library(
        name = name,
        srcs = [name_pb],
        deps = [
            str(Label("//android:grpc_deps")),
        ],
        exports = [
            str(Label("//android:grpc_deps")),
        ],
        visibility = visibility,
    )
