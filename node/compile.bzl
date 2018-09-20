load("//:compile.bzl", "proto_compile")

def node_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//node:js"))],
        **kwargs
    )

def grpc_node_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//node:js")), str(Label("//node:grpc_js"))],
        **kwargs
    )
