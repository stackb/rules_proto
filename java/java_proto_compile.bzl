load("//:compile.bzl", "proto_compile")

def java_proto_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//java:java")),
        ],
        **kwargs
    )
