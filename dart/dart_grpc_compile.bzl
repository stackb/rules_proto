load("//:compile.bzl", "proto_compile")

def dart_grpc_compile(**kwargs):
    # Append the dart plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//dart:grpc_dart"),
    ]
    proto_compile(**kwargs)
