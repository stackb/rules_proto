load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def proto_deps():
    # todo: put protobuf here
    pass

def grpc_deps():
    http_archive(
        name = "com_github_grpc_grpc",
        strip_prefix = "grpc-1.15.0",
        url = "https://github.com/grpc/grpc/archive/v1.15.0.tar.gz",
        sha256 = "013cc34f3c51c0f87e059a12ea203087a7a15dca2e453295345e1d02e2b9634b",
    )
