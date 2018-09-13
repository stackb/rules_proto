load("@//:rules.bzl", "proto_compile")

def go_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//go:go"],
        **kwargs
    )

def grpc_go_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//go:grpc_go"],
        **kwargs
    )

