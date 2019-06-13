load("//:compile.bzl", "proto_compile")


def ruby_grpc_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//ruby:ruby")),
            str(Label("//ruby:grpc_ruby")),
        ],
        **kwargs
    )
