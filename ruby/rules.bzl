load("@//:rules.bzl", "proto_compile")

def ruby_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//ruby:ruby"],
        grpc_plugins = ["//ruby:grpc_ruby"],
        **kwargs
    )
