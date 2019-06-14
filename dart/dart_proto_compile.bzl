load("//:compile.bzl", "proto_compile")

def dart_proto_compile(**kwargs):
    # Append the dart plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//dart:dart"),
    ]
    proto_compile(**kwargs)
