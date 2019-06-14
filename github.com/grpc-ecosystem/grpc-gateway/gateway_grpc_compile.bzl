load("//:compile.bzl", "proto_compile")

def gateway_grpc_compile(**kwargs):
    # Append the grpc-gateway plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//github.com/grpc-ecosystem/grpc-gateway:grpc-gateway"),
    ]
    proto_compile(**kwargs)
