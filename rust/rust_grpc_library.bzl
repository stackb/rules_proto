load("//rust:rust_grpc_compile.bzl", "rust_grpc_compile")
load("//rust:rust_proto_lib.bzl", "rust_proto_lib")
load("@io_bazel_rules_rust//rust:rust.bzl", "rust_library")

def rust_grpc_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    name_lib = kwargs.get("name") + "_lib"
    rust_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create lib file
    rust_proto_lib(
        name = name_lib,
        compilation = name_pb,
        grpc = True,
    )

    # Create rust library
    rust_library(
        name = kwargs.get("name"),
        srcs = [name_pb, name_lib],
        deps = [
            "@io_bazel_rules_rust//proto/raze:protobuf",
            "@io_bazel_rules_rust//proto/raze:grpc",
            "@io_bazel_rules_rust//proto/raze:tls_api",
            "@io_bazel_rules_rust//proto/raze:tls_api_stub",
        ],
        visibility = kwargs.get("visibility"),
    )
