load("//:compile.bzl", "proto_compile")

def objc_proto_compile(**kwargs):
    # Prepend the objc plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//objc:objc"),
    ]
    proto_compile(**kwargs)
