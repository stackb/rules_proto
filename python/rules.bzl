load("@//:rules.bzl", "proto_compile")

def py_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//python:python"],
        grpc_plugins = ["//python:grpc_python"],
        **kwargs
    )
