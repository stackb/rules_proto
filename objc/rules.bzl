load("@//:rules.bzl", "proto_compile")

def objc_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//objc:objc"],
        grpc_plugins = ["//objc:grpc_objc"],
        **kwargs
    )
