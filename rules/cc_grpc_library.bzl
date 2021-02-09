load(
    "@build_stack_rules_proto//rules:cc_proto_compile.bzl",
    "cc_proto_compile",
)
load("@rules_cc//cc:defs.bzl", "cc_library")

GRPC_DEPS = [
    "@com_google_protobuf//:protoc_lib",
    "@com_github_grpc_grpc//:grpc++",
    "@com_github_grpc_grpc//:grpc++_reflection",
]

def cc_grpc_library(**kwargs):
    name_pb = kwargs.get("name") + "_pb"

    cc_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    cc_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = GRPC_DEPS,
        includes = [name_pb],
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags", []),
    )
