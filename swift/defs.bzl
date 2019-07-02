# Aggregate all `swift` rules to one loadable file
load(":swift_proto_compile.bzl", "swift_proto_compile")
load(":swift_grpc_compile.bzl", "swift_grpc_compile")
load(":swift_proto_library.bzl", "swift_proto_library")
load(":swift_grpc_library.bzl", "swift_grpc_library")
