# Aggregate all `dart` rules to one loadable file
load(":dart_proto_compile.bzl", "dart_proto_compile")
load(":dart_grpc_compile.bzl", "dart_grpc_compile")
load(":dart_proto_library.bzl", "dart_proto_library")
load(":dart_grpc_library.bzl", "dart_grpc_library")
