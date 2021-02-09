load(
    "@build_stack_rules_proto//rules:py_grpc_compile.bzl",
    "py_grpc_compile",
)
load("@rules_python//python:defs.bzl", "py_library")

GRPC_DEPS = [
    "@com_google_protobuf//:protobuf_python",
    "@com_github_grpc_grpc//src/python/grpcio/grpc:grpcio",
]

def py_grpc_library(**kwargs):
    name_pb = kwargs.get("name") + "_pb"

    py_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    py_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = GRPC_DEPS,
        imports = [name_pb],
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags", []),
    )
