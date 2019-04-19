load("//dart:dart_grpc_compile.bzl", "dart_grpc_compile")
load("@io_bazel_rules_dart//dart/build_rules:core.bzl", "dart_library")

def dart_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    dart_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

	dart_library(
        name = name,
        srcs = [name_pb],
        deps = [
            str(Label("@vendor_protobuf//:protobuf")),
            str(Label("@vendor_grpc//:grpc")),
        ],
        pub_pkg_name = name,
        visibility = visibility,
    )
