load(
    "@build_stack_rules_proto//rules:java_proto_compile.bzl",
    "java_proto_compile",
)
load("@rules_java//java:defs.bzl", "java_library")

PROTO_DEPS = ["@com_google_protobuf//:protobuf_java"]

def java_proto_library(**kwargs):
    name_pb = kwargs.get("name") + "_pb"

    java_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    java_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = PROTO_DEPS,
        exports = PROTO_DEPS,
        visibility = kwargs.get("visibility"),
        tags = kwargs.get("tags", []),
    )
