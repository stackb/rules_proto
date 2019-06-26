load("//python:python_grpc_compile.bzl", "python_grpc_compile")
load("@protobuf_py_deps//:requirements.bzl", protobuf_requirements = "all_requirements")
load("@grpc_py_deps//:requirements.bzl", grpc_requirements = "all_requirements")

def python_grpc_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    python_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k != "name"} # Forward args except name
    )

    # Create python library
    native.py_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = depset(protobuf_requirements + grpc_requirements).to_list(),
        imports = [name_pb],
        visibility = kwargs.get("visibility"),
    )

# Alias
py_grpc_library = python_grpc_library
