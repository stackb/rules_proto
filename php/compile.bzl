load("@//:compile.bzl", "proto_compile")

def php_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//php:php"],
        **kwargs
    )

def grpc_php_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//php:php", "//php:grpc_php"],
        **kwargs
    )
