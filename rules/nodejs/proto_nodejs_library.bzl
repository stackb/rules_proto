"proto_nodejs_library.bzl provides a js_library for proto files."

load("@build_bazel_rules_nodejs//:index.bzl", "js_library")

def proto_nodejs_library(**kwargs):
    js_library(**kwargs)
