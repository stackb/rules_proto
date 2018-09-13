load("@//:rules.bzl", "proto_compile")

def cpp_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//cpp:cpp"],
        **kwargs
    )

def grpc_cpp_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//cpp:cpp", "//cpp:grpc_cpp"],
        **kwargs
    )
