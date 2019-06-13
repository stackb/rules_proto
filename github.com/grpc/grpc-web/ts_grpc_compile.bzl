load("//:compile.bzl", "proto_compile")

def ts_grpc_compile(**kwargs):
    # Prepend the grpc-web plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//github.com/grpc/grpc-web:ts"),
    ]
    proto_compile(**kwargs)
