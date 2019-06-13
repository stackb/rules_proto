load("//:compile.bzl", "proto_compile")

def swift_grpc_compile(**kwargs):
    # Prepend the swift plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//swift:grpc_swift"),
    ]
    proto_compile(**kwargs)
