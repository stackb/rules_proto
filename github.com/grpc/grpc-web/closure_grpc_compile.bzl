load("//:compile.bzl", "proto_compile")

def closure_grpc_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//github.com/grpc/grpc-web:closure")),
        ],
        **kwargs
    )
