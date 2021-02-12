load(
    "@build_stack_rules_proto//rules:nodejs_proto_compile.bzl",
    "nodejs_proto_compile",
)
load(
    "@build_stack_rules_proto//rules:proto_compile_js_library.bzl",
    "proto_compile_js_library",
)

PROTO_DEPS = [
    "@google_protobuf_node_modules//google-protobuf",
]

def nodejs_proto_library(**kwargs):
    name_pb = kwargs.get("name") + "_pb"

    nodejs_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")}  # Forward args
    )

    proto_compile_js_library(
        name = kwargs.get("name"),
        deps = [name_pb],
        # js_deps = PROTO_DEPS,
        visibility = kwargs.get("visibility", []),
        tags = kwargs.get("tags", []),
    )
