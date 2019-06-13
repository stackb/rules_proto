load("//:compile.bzl", "proto_compile")

def python_proto_compile(**kwargs):
    # Prepend the python plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//python:python"),
    ]
    proto_compile(**kwargs)
