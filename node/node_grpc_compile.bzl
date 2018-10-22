load("//:compile.bzl", "proto_compile")

def node_grpc_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//node:js")),
            str(Label("//node:grpc_js")),
        ],
        **kwargs
    )
