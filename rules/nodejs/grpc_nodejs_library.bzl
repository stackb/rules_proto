"grpc_nodejs_library.bzl provides a js_library for grpc files."

load("@build_bazel_rules_nodejs//:index.bzl", "js_library")

def grpc_nodejs_library(**kwargs):
    js_library(**kwargs)
