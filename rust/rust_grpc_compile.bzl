load("//:compile.bzl", "proto_compile")

def rust_grpc_compile(**kwargs):
    # Append the rust plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//rust:rust"),
        Label("//rust:grpc_rust"),
    ]
    proto_compile(**kwargs)
