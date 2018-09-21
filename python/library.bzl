load("//python:compile.bzl", "py_proto_compile", "py_grpc_compile")
load("@grpc_py_deps//:requirements.bzl", "all_requirements")

def py_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    py_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = verbose,
    )
    native.py_library(
        name = name,
        srcs = [name_pb],
        deps = all_requirements, # fixme don't need grpc here
        # This magically adds REPOSITORY_NAME/PACKAGE_NAME/{name_pb} to PYTHONPATH
        imports = [name_pb],
        visibility = visibility,
    )

def py_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    py_grpc_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        verbose = verbose,
    )
    native.py_library(
        name = name,
        srcs = [name_pb],
        deps = all_requirements,
        # This magically adds REPOSITORY_NAME/PACKAGE_NAME/{name_pb} to PYTHONPATH
        imports = [name_pb],
        visibility = visibility,
    )