load("//:compile.bzl", "proto_compile")

def closure_grpc_compile(**kwargs):
    # Append the grpc.js plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//github.com/stackb/grpc.js:grpc.js"),
    ]
    proto_compile(**kwargs)
