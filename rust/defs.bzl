# Aggregate all `rust` rules to one loadable file
load(":rust_proto_compile.bzl", "rust_proto_compile")
load(":rust_grpc_compile.bzl", "rust_grpc_compile")
load(":rust_proto_library.bzl", "rust_proto_library")
load(":rust_grpc_library.bzl", "rust_grpc_library")
