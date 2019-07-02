# Aggregate all `csharp` rules to one loadable file
load(":csharp_proto_compile.bzl", "csharp_proto_compile")
load(":csharp_grpc_compile.bzl", "csharp_grpc_compile")
load(":csharp_proto_library.bzl", "csharp_proto_library")
load(":csharp_grpc_library.bzl", "csharp_grpc_library")
