load("@//:rules.bzl", "proto_compile")

def node_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//node:js"],
        **kwargs
    )

def grpc_node_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//node:js", "//node:grpc_js"],
        **kwargs
    )
