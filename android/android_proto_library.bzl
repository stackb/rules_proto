load("//android:android_proto_compile.bzl", "android_proto_compile")
load("@build_bazel_rules_android//android:rules.bzl", "android_library")

def android_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"

    android_proto_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )

    android_library(
        name = name,
        srcs = [name_pb],
        deps = [
            str(Label("//android:proto_deps")),
        ],
        exports = [
            str(Label("//android:proto_deps")),
        ],
        visibility = visibility,
    )

