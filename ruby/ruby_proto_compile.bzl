load("//:compile.bzl", "proto_compile")


def ruby_proto_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//ruby:ruby")),
        ],
        **kwargs
    )
