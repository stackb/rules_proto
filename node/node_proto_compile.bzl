load("//:compile.bzl", "proto_compile")

def node_proto_compile(**kwargs):
    # Append the node plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//node:js"),
    ]
    proto_compile(**kwargs)
