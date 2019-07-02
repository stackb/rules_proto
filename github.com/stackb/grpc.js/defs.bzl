# Aggregate all `grpc.js` rules to one loadable file
load(":grpcjs_grpc_compile.bzl", _grpcjs_grpc_compile="grpcjs_grpc_compile")
load(":grpcjs_grpc_library.bzl", _grpcjs_grpc_library="grpcjs_grpc_library")

grpcjs_grpc_compile = _grpcjs_grpc_compile
grpcjs_grpc_library = _grpcjs_grpc_library
