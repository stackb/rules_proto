# Aggregate all `scala` rules to one loadable file
load(":scala_proto_compile.bzl", "scala_proto_compile")
load(":scala_grpc_compile.bzl", "scala_grpc_compile")
load(":scala_proto_library.bzl", "scala_proto_library")
load(":scala_grpc_library.bzl", "scala_grpc_library")
