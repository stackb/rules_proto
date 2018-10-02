load("//:compile.bzl", "proto_compile")

def scala_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//scala:scala"))],
        **kwargs
    )

def scala_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//scala:grpc_scala"))],
        **kwargs
    )
