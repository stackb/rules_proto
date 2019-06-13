load("//:compile.bzl", "proto_compile")

def objc_grpc_compile(**kwargs):
    # Prepend the objc plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//objc:objc"),
        Label("//objc:grpc_objc"),
    ]
    proto_compile(**kwargs)
