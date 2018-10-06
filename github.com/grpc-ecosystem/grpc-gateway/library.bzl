load("//github.com/grpc-ecosystem/grpc-gateway:compile.bzl", "grpc_gateway_compile")
load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@io_bazel_rules_go//proto:def.bzl", "go_proto_library")

def grpc_gateway_proto_library(**kwargs):
    name = kwargs.get("name")
    importpath = kwargs.get("importpath")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

    compilers = kwargs.get("compilers")
    if not compilers:
        compilers = [
            "@io_bazel_rules_go//proto:go_grpc",
            "@com_github_grpc_ecosystem_grpc_gateway//protoc-gen-grpc-gateway:go_gen_grpc_gateway",
        ]

    go_proto_library(
        name = name,
        compilers = compilers,
        importpath = importpath,
        proto = deps[0],
        visibility = visibility,
    )

def grpc_gateway_library(**kwargs):
    name = kwargs.get("name")
    importpath = kwargs.get("importpath")
    visibility = kwargs.get("visibility")

    name_go = name + "_go"
    name_gw = name + "_gw"

    kwargs["name"] = name_go
    grpc_gateway_proto_library(**kwargs)

    go_library(
        name = name,
        embed = [name_go],
        deps = [
            "@org_golang_google_genproto//googleapis/api:go_default_library",
            "@org_golang_google_genproto//googleapis/api/annotations:go_default_library",
        ],
        importpath = importpath,
        visibility = visibility,
    )
