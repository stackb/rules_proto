load("//python:python_grpc_compile.bzl", "python_grpc_compile")
load("@grpc_py_deps//:requirements.bzl", grpc_requirements = "all_requirements")

def python_grpc_library(**kwargs):
    # Compile protos
    name_pb = kwargs.get("name") + "_pb"
    python_grpc_compile(
        name = name_pb,
        **{k: v for (k, v) in kwargs.items() if k in ("deps", "verbose")} # Forward args
    )

    # Create python library
    native.py_library(
        name = kwargs.get("name"),
        srcs = [name_pb],
        deps = [
            "@com_google_protobuf//:protobuf_python",
        ] + grpc_requirements,
        imports = [name_pb],
        visibility = kwargs.get("visibility"),
    )

# Alias
py_grpc_library = python_grpc_library
