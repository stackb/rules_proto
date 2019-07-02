# Aggregate all `rust` rules to one loadable file
load(":rust_proto_compile.bzl", _rust_proto_compile="rust_proto_compile")
load(":rust_grpc_compile.bzl", _rust_grpc_compile="rust_grpc_compile")
load(":rust_proto_library.bzl", _rust_proto_library="rust_proto_library")
load(":rust_grpc_library.bzl", _rust_grpc_library="rust_grpc_library")

rust_proto_compile = _rust_proto_compile
rust_grpc_compile = _rust_grpc_compile
rust_proto_library = _rust_proto_library
rust_grpc_library = _rust_grpc_library
