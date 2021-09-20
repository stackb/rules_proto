"grpc_closure_js_library.bzl provides a closure_js_library for grpc files."

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def grpc_closure_js_library(**kwargs):
    closure_js_library(**kwargs)
