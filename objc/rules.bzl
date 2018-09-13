load("@//:rules.bzl", "proto_compile")

def objc_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//objc:objc"],
        **kwargs
    )

def grpc_objc_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//objc:objc", "//objc:grpc_objc"],
        **kwargs
    )
