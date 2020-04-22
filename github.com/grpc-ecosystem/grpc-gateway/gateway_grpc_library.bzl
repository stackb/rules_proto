load("//github.com/grpc-ecosystem/grpc-gateway:gateway_grpc_compile.bzl", "gateway_grpc_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

def gateway_grpc_library(**kwargs):
    name = kwargs.get("name")
    importpath = kwargs.get("importpath")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    compilers = kwargs.get("compilers")
    if not compilers:
        compilers = [
            "@io_bazel_rules_go//proto:go_grpc",
            "@grpc_ecosystem_grpc_gateway//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
        ]

    go_proto_library(
        name = name,
        compilers = compilers,
        importpath = importpath,
        proto = deps[0],
        deps = ["@go_googleapis//google/api:annotations_go_proto"] + deps[1:],
        visibility = visibility,
    )
