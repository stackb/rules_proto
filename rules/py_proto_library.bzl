load(
    "@build_stack_rules_proto//rules:py_proto_compile.bzl",
    "py_proto_compile",
)
load("@rules_python//python:defs.bzl", "py_library")

PROTO_DEPS = ["@com_google_protobuf//:protobuf_python"]

def py_proto_library(**kwargs):
    name_pb = kwargs.get("name") + "_pb"

    py_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    py_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = PROTO_DEPS,
        imports = [name_pb],
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags", []),
    )
