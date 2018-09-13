load("@//:rules.bzl", "proto_compile")

def cpp_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//cpp:cpp"],
        grpc_plugins = ["//cpp:grpc_cpp"],
        **kwargs
    )
