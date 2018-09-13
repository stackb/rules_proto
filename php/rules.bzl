load("@//:rules.bzl", "proto_compile")

def php_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//php:php"],
        grpc_plugins = ["//php:grpc_php"],
        **kwargs
    )
