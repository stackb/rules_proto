load("//github.com/grpc-ecosystem/grpc-gateway:gateway_grpc_compile.bzl", "gateway_grpc_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

def gateway_grpc_library(**kwargs):
    # Apply default args
    if not kwargs.get("compilers"):
        kwargs["compilers"] = [
            "@io_bazel_rules_go//proto:go_grpc",
            "@grpc_ecosystem_grpc_gateway//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
        ]

    # Create go library
    go_proto_library(
        proto = kwargs.get("deps")[0],
        deps = ["@go_googleapis//google/api:annotations_go_proto"],
        **{k: v for (k, v) in kwargs.items() if k != "deps"} # Forward args except deps
    )
