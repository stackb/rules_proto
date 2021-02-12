load(
    "@build_stack_rules_proto//rules:nodejs_proto_compile.bzl",
    "nodejs_proto_compile",
)
load("@build_bazel_rules_nodejs//:index.bzl", "js_library")

PROTO_DEPS = [
    "@google_protobuf_node_modules//google-protobuf",
]

def nodejs_proto_library(**kwargs):
    name_pb = kwargs.get("name") + "_pb"
    js_deps = kwargs.pop("js_deps", PROTO_DEPS)

    nodejs_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    js_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = js_deps,
        visibility = kwargs.get("visibility", []),
        package_name = kwargs.get("name"),
        tags = kwargs.get("tags", []),
    )
