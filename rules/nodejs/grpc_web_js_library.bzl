"grpc_web_js_library.bzl provides a js_library for grpc files."

load("@build_bazel_rules_nodejs//:index.bzl", "js_library")

def grpc_web_js_library(**kwargs):
    js_library(**kwargs)
