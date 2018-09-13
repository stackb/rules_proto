load("@//:rules.bzl", "proto_compile")

def csharp_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//csharp:csharp"],
        grpc_plugins = ["//csharp:grpc_csharp"],
        **kwargs
    )
