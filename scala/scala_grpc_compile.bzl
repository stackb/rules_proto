load("//:compile.bzl", "proto_compile")

def scala_grpc_compile(**kwargs):
    # Append the scala plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//scala:grpc_scala"),
    ]
    proto_compile(**kwargs)
