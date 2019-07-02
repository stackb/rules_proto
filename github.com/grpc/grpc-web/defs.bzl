# Aggregate all `grpc-web` rules to one loadable file
load(":closure_grpc_compile.bzl", "closure_grpc_compile")
load(":commonjs_grpc_compile.bzl", "commonjs_grpc_compile")
load(":commonjs_dts_grpc_compile.bzl", "commonjs_dts_grpc_compile")
load(":ts_grpc_compile.bzl", "ts_grpc_compile")
load(":closure_grpc_library.bzl", "closure_grpc_library")
