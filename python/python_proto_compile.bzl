load("//:compile.bzl", "proto_compile")

def python_proto_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//python:python")),
        ],
        **kwargs
    )
