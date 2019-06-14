load("//:compile.bzl", "proto_compile")

def scala_proto_compile(**kwargs):
    # Append the scala plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//scala:scala"),
    ]
    proto_compile(**kwargs)
