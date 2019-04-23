load("//rust:rust_proto_compile.bzl", "rust_proto_compile")
load("//rust:rust_proto_lib.bzl", "rust_proto_lib")
load("@io_bazel_rules_rust//rust:rust.bzl", "rust_library")

def rust_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    name_lib = name + "_lib"

    rust_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = kwargs.pop("verbose", 0),
        transitivity = kwargs.pop("transitivity", {}),
        transitive = kwargs.pop("transitive", True),
    )

    rust_proto_lib(
        name = name_lib,
        compilation = name_pb,
    )

    rust_library(
        name = name,
        srcs = [name_pb, name_lib],
        deps = [
            "@io_bazel_rules_rust//proto/raze:protobuf",
        ],
        visibility = visibility,
    )
