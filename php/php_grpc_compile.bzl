load("//:compile.bzl", "proto_compile")

def php_grpc_compile(**kwargs):
    # Append the php plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//php:php"),
        Label("//php:grpc_php"),
    ]
    proto_compile(**kwargs)
