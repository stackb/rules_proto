load("//:compile.bzl", "proto_compile")

def grpc_gateway_compile(**kwargs):
    plugins = kwargs.get("plugins")
    if not plugins:
        kwargs["plugins"] = [str(Label("//github.com/grpc-ecosystem/grpc-gateway:grpc-gateway"))]
    proto_compile(**kwargs)


def grpc_gateway_swagger_compile(**kwargs):
    plugins = kwargs.get("plugins")
    if not plugins:
        kwargs["plugins"] = [str(Label("//github.com/grpc-ecosystem/grpc-gateway:swagger"))]
    proto_compile(**kwargs)
