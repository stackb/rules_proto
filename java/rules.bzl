load("@//:rules.bzl", "proto_compile")

def java_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//java:java"],
        grpc_plugins = ["//java:grpc_java"],
        **kwargs
    )
