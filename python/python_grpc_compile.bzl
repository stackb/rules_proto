load("//:compile.bzl", "proto_compile")


def python_grpc_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//python:python")),
            str(Label("//python:grpc_python")),
        ],
        **kwargs
    )