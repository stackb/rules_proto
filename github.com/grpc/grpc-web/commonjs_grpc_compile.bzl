load("//:compile.bzl", "proto_compile")

def commonjs_grpc_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//github.com/grpc/grpc-web:commonjs")),
        ],
        **kwargs
    )
