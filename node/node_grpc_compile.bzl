load("//:compile.bzl", "proto_compile")

def node_grpc_compile(**kwargs):
    # Prepend the node plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//node:js"),
        Label("//node:grpc_js"),
    ]
    proto_compile(**kwargs)
