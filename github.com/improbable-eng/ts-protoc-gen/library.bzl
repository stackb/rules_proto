load("//node:compile.bzl", "node_proto_compile", "node_grpc_compile")
load("//github.com/improbable-eng/ts-protoc-gen:compile.bzl", "ts_proto_compile", "ts_grpc_compile")
load("@build_bazel_rules_typescript//:defs.bzl", "ts_library")

def ts_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    name_js = name + "_js"

    node_proto_compile(
        name = name_js,
        deps = deps,
        visibility = visibility,
    )
    
    ts_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = verbose,
    )

    ts_library(
        name = name,
        srcs = [name_pb],
        deps = [
        ],
        visibility = visibility,
    )

def ts_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")

    name_js = name + "_pbjs"
    name_pb_grpc = name + "_pb_grpc"

    node_grpc_compile(
        name = name_js,
        deps = deps,
        visibility = visibility,
    )
    
    ts_grpc_compile(
        name = name_pb_grpc,
        deps = deps,
        visibility = visibility,
        verbose = verbose,
    )

    ts_library(
        name = name,
        srcs = [name_pb_grpc],
        deps = [
        ],
        visibility = visibility,
    )