load("//:compile.bzl", "proto_compile")

def php_proto_compile(**kwargs):
    # Append the php plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//php:php"),
    ]
    proto_compile(**kwargs)
