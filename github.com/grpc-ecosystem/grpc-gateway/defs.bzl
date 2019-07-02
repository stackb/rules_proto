# Aggregate all `grpc-gateway` rules to one loadable file
load(":gateway_grpc_compile.bzl", "gateway_grpc_compile")
load(":gateway_swagger_compile.bzl", "gateway_swagger_compile")
load(":gateway_grpc_library.bzl", "gateway_grpc_library")
