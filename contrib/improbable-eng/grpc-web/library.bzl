load("//node:compile.bzl", "node_proto_compile")
load("//ts:compile.bzl", "ts_proto_compile", "grpc_ts_proto_compile")
load("@build_bazel_rules_typescript//:defs.bzl", "ts_library")

def grpc_web_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    name_pb_grpc = name + "_pb_grpc"

    grpc_node_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
    )
    
    grpc_ts_proto_compile(
        name = name_pb_grpc,
        deps = deps,
        visibility = visibility,
        verbose = verbose,
    )

    ts_library(
        name = name,
        srcs = [name_pb, name_pb_grpc],
        deps = [
        ],
        visibility = visibility,
    )