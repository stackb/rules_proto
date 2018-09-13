load("@//:rules.bzl", "proto_compile")

def py_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//python:python"],
        **kwargs
    )

def grpc_py_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//python:python", "//python:grpc_python"],
        **kwargs
    )
