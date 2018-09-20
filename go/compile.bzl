load("//:compile.bzl", "proto_compile")

def go_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//go:go"))],
        **kwargs
    )

def grpc_go_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//go:grpc_go"))],
        **kwargs
    )
