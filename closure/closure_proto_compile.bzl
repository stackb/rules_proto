load("//:compile.bzl", "proto_compile")

def closure_proto_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//closure:js")),
        ],
        **kwargs
    )
