load("//:compile.bzl", "proto_compile")

def grpc_gateway_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/grpc-ecosystem/grpc-gateway:grpc-gateway"))],
        **kwargs
    )

def grpc_gateway_swagger_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/grpc-ecosystem/grpc-gateway:swagger"))],
        **kwargs
    )
