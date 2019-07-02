# Aggregate all `grpc-gateway` rules to one loadable file
load(":gateway_grpc_compile.bzl", _gateway_grpc_compile="gateway_grpc_compile")
load(":gateway_swagger_compile.bzl", _gateway_swagger_compile="gateway_swagger_compile")
load(":gateway_grpc_library.bzl", _gateway_grpc_library="gateway_grpc_library")

gateway_grpc_compile = _gateway_grpc_compile
gateway_swagger_compile = _gateway_swagger_compile
gateway_grpc_library = _gateway_grpc_library
