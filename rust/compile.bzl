load("//:compile.bzl", "proto_compile")

def rust_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//rust:rust"))],
        **kwargs
    )

def rust_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//rust:rust")), str(Label("//rust:grpc_rust"))],
        **kwargs
    )
