load("//:compile.bzl", "proto_compile")

def gateway_swagger_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//github.com/grpc-ecosystem/grpc-gateway:swagger")),
        ],
        **kwargs
    )
