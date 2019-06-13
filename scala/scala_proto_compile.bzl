load("//:compile.bzl", "proto_compile")

def scala_proto_compile(**kwargs):
    # Prepend the scala plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//scala:scala"),
    ]
    proto_compile(**kwargs)
