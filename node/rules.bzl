load("@//:rules.bzl", "proto_compile")

def node_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//node:js"],
        grpc_plugins = ["//node:grpc_js"],
        **kwargs
    )
