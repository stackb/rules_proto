load("@//:compile.bzl", "proto_compile")

def java_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//java:java"],
        **kwargs
    )

def grpc_java_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//java:java", "//java:grpc_java"],
        **kwargs
    )

