load("//:compile.bzl", "proto_compile")

def rust_grpc_compile(**kwargs):
    # Prepend the rust plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//rust:rust"),
        Label("//rust:grpc_rust"),
    ]
    proto_compile(**kwargs)
