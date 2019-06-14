load("//:compile.bzl", "proto_compile")

def ruby_proto_compile(**kwargs):
    # Append the ruby plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//ruby:ruby"),
    ]
    proto_compile(**kwargs)
