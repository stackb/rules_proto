load("//:compile.bzl", "proto_compile")

def python_proto_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//python:python"))],
        **kwargs
    )

def python_grpc_compile(**kwargs):
    proto_compile(
        plugins = [str(Label("//python:python")), str(Label("//python:grpc_python"))],
        **kwargs
    )

# Alias to shorter names 
py_proto_compile = python_proto_compile
py_grpc_compile = python_grpc_compile
