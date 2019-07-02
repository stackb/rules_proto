# Aggregate all `python` rules to one loadable file
load(":python_proto_compile.bzl", "python_proto_compile")
load(":python_grpc_compile.bzl", "python_grpc_compile")
load(":python_proto_library.bzl", "python_proto_library")
load(":python_grpc_library.bzl", "python_grpc_library")
