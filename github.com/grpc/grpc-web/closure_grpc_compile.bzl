load("//:compile.bzl", "proto_compile")

def closure_grpc_compile(**kwargs):
    # Append the grpc-web plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//github.com/grpc/grpc-web:closure"),
    ]
    proto_compile(**kwargs)
