load("//android:android_grpc_compile.bzl", "android_grpc_compile")
load("@build_bazel_rules_android//android:rules.bzl", "android_library")

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
        deps = GRPC_DEPS,
        exports = GRPC_DEPS,
        visibility = kwargs.get("visibility"),
    )

GRPC_DEPS = [
    "@com_google_guava_guava_android//jar",
    "@com_google_protobuf//:protobuf_javalite",
    "@com_google_protobuf//:protobuf_java_util",
    "@javax_annotation_javax_annotation_api//jar",
    "@io_grpc_grpc_java//core",
    "@io_grpc_grpc_java//protobuf-lite",
    "@io_grpc_grpc_java//stub",
]
