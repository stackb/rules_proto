load("@//:compile.bzl", "proto_compile")

def py_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//python:python"],
        **kwargs
    )

def grpc_py_proto_compile(**kwargs):
    proto_compile(
        plugins = ["//python:python", "//python:grpc_python"],
        **kwargs
    )


def grpc_py_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    grpc_py_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = verbose,
    )
    native.py_library(
        name = name,
        srcs = [name_pb],
        deps = ["@//python:grpc_deps"],
        # This magically adds REPOSITORY_NAME/PACKAGE_NAME/{name_pb} to PYTHONPATH
        imports = [name_pb],
        visibility = visibility,
    )

