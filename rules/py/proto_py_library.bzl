"proto_py_library.bzl provides a py_library for proto files."

load("@rules_python//python:defs.bzl", "py_library")

def proto_py_library(**kwargs):
    py_library(**kwargs)
