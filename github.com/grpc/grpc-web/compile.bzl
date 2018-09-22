load("//:compile.bzl", "proto_compile")

def grpc_web_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//github.com/grpc/grpc-web:grpc-web"))],
        **kwargs
    )
