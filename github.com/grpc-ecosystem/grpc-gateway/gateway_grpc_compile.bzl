load("//:compile.bzl", "proto_compile")

def gateway_grpc_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//github.com/grpc-ecosystem/grpc-gateway:grpc-gateway")),
        ],
        **kwargs
    )
