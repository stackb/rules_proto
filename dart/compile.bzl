load("@//:compile.bzl", "proto_compile")

def dart_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//dart:dart"],
        **kwargs
    )

def grpc_dart_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//dart:grpc_dart"],
        **kwargs
    )


