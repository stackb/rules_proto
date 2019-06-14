load("//:compile.bzl", "proto_compile")

def closure_proto_compile(**kwargs):
    # Append the closure plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//closure:js"),
    ]
    proto_compile(**kwargs)
