load("@//:rules.bzl", "proto_compile")

def py_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//py:python"],
        grpc_plugins = ["//py:grpc_python"],
        **kwargs
    )
