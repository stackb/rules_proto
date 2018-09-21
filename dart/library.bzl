load("//dart:compile.bzl", "dart_grpc_compile")
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
        verbose = verbose,
    )
    dart_library(
        name = name,
        srcs = [name_pb],
        deps = [
            "@vendor_protobuf//:protobuf",
        ],
        lib_root = ".",
        pub_pkg_name = "foo",
        visibility = visibility,
    )

