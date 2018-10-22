load("//:compile.bzl", "proto_compile")

def dart_proto_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//dart:dart")),
        ],
        **kwargs
    )
