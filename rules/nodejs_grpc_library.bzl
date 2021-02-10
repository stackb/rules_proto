load(
    "@build_stack_rules_proto//rules:nodejs_grpc_compile.bzl",
    "nodejs_grpc_compile",
)

load("@build_bazel_rules_nodejs//:index.bzl", "js_library")

GRPC_DEPS = [
    "@nodejs_proto_grpc_modules//google-protobuf",
    "@nodejs_proto_grpc_modules//@grpc/grpc-js",
]

def nodejs_grpc_library(**kwargs):
    name_pb = kwargs.get("name") + "_pb"

    nodejs_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    js_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = GRPC_DEPS,
        package_name = kwargs.get("name"),
        visibility = kwargs.get("visibility", []),
        tags = kwargs.get("tags", []),
    )
