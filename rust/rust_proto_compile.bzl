load("//:compile.bzl", "proto_compile")

def rust_proto_compile(**kwargs):
    # Append the rust plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//rust:rust"),
    ]
    proto_compile(**kwargs)
