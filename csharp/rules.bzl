load("@//:rules.bzl", "proto_compile")

def csharp_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//csharp:csharp"],
        **kwargs
    )

def grpc_csharp_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//csharp:csharp", "//csharp:grpc_csharp"],
        **kwargs
    )
