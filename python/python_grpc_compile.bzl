load("//:compile.bzl", "proto_compile")

def python_grpc_compile(**kwargs):
    # Append the python plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//python:python"),
        Label("//python:grpc_python"),
    ]
    proto_compile(**kwargs)
