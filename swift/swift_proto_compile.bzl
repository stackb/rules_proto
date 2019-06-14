load("//:compile.bzl", "proto_compile")

def swift_proto_compile(**kwargs):
    # Append the swift plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//swift:swift"),
    ]
    proto_compile(**kwargs)
