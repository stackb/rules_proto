load(
    "@build_stack_rules_proto//rules:nodejs_proto_compile.bzl",
    "nodejs_proto_compile",
)
load("@build_bazel_rules_nodejs//:index.bzl", "js_library")

PROTO_DEPS = [
    "@nodejs_proto_grpc_modules//google-protobuf",
]

def nodejs_proto_library(**kwargs):
    name_pb = kwargs.get("name") + "_pb"

    nodejs_proto_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")}  # Forward args
    )

    js_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = PROTO_DEPS,
        package_name = kwargs.get("name"),
        visibility = kwargs.get("visibility", []),
        tags = kwargs.get("tags", []),
    )
