load("@//:rules.bzl", "proto_compile")

def closure_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//closure:js"],
        **kwargs
    )
