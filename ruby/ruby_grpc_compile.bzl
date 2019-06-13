load("//:compile.bzl", "proto_compile")

def ruby_grpc_compile(**kwargs):
    # Prepend the ruby plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//ruby:ruby"),
        Label("//ruby:grpc_ruby"),
    ]
    proto_compile(**kwargs)
