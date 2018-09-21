load("//:compile.bzl", "proto_compile")

def ruby_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//ruby:ruby"))],
        **kwargs
    )

def grpc_ruby_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//ruby:ruby")), str(Label("//ruby:grpc_ruby"))],
        **kwargs
    )
