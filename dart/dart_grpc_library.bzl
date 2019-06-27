load("//dart:dart_grpc_compile.bzl", "dart_grpc_compile")
load("@io_bazel_rules_dart//dart/build_rules:core.bzl", "dart_library")

def dart_grpc_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    dart_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create dart library
    dart_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            Label("@vendor_protobuf//:protobuf"),
            Label("@vendor_grpc//:grpc"),
        ],
        pub_pkg_name = kwargs.get("name"),
        visibility = kwargs.get("visibility"),
    )
