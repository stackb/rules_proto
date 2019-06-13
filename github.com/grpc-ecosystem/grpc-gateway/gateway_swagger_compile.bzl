load("//:compile.bzl", "proto_compile")

def gateway_swagger_compile(**kwargs):
    # Prepend the grpc-gateway plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//github.com/grpc-ecosystem/grpc-gateway:swagger"),
    ]
    proto_compile(**kwargs)
