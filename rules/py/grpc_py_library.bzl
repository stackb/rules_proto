"grpc_py_library.bzl provides a py_library for grpc files."

load("@rules_python//python:defs.bzl", "py_library")

def grpc_py_library(**kwargs):
    py_library(**kwargs)
