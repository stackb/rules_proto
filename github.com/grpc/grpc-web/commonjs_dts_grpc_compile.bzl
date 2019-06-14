load("//:compile.bzl", "proto_compile")

def commonjs_dts_grpc_compile(**kwargs):
    # Append the grpc-web plugins and call generic compile
    kwargs["plugins"] = kwargs.get("plugins", []) + [
        Label("//github.com/grpc/grpc-web:commonjs_dts"),
    ]
    proto_compile(**kwargs)
