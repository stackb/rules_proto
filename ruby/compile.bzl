load("@//:compile.bzl", "proto_compile")

def ruby_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//ruby:ruby"],
        **kwargs
    )

def grpc_ruby_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//ruby:ruby", "//ruby:grpc_ruby"],
        **kwargs
    )
