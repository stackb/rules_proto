load("//:compile.bzl", "proto_compile")

def php_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//php:php"))],
        **kwargs
    )

def grpc_php_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//php:php")), str(Label("//php:grpc_php"))],
        **kwargs
    )
