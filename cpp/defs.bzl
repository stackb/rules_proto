# Aggregate all `cpp` rules to one loadable file
load(":cpp_proto_compile.bzl", _cpp_proto_compile="cpp_proto_compile")
load(":cpp_grpc_compile.bzl", _cpp_grpc_compile="cpp_grpc_compile")
load(":cpp_proto_library.bzl", _cpp_proto_library="cpp_proto_library")
load(":cpp_grpc_library.bzl", _cpp_grpc_library="cpp_grpc_library")

cpp_proto_compile = _cpp_proto_compile
cpp_grpc_compile = _cpp_grpc_compile
cpp_proto_library = _cpp_proto_library
cpp_grpc_library = _cpp_grpc_library
