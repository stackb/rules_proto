load("//android:android_grpc_compile.bzl", "android_grpc_compile")
load("@build_bazel_rules_android//android:android_library.bzl", "android_library")

def android_grpc_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    android_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create android library
    android_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            Label("//android:grpc_deps"),
        ],
        exports = [
            Label("//android:grpc_deps"),
        ],
        visibility = kwargs.get("visibility"),
    )
