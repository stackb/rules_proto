# Aggregate all `csharp` rules to one loadable file
load(":csharp_proto_compile.bzl", _csharp_proto_compile="csharp_proto_compile")
load(":csharp_grpc_compile.bzl", _csharp_grpc_compile="csharp_grpc_compile")
load(":csharp_proto_library.bzl", _csharp_proto_library="csharp_proto_library")
load(":csharp_grpc_library.bzl", _csharp_grpc_library="csharp_grpc_library")

csharp_proto_compile = _csharp_proto_compile
csharp_grpc_compile = _csharp_grpc_compile
csharp_proto_library = _csharp_proto_library
csharp_grpc_library = _csharp_grpc_library
