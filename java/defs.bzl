# Aggregate all `java` rules to one loadable file
load(":java_proto_compile.bzl", "java_proto_compile")
load(":java_grpc_compile.bzl", "java_grpc_compile")
load(":java_proto_library.bzl", "java_proto_library")
load(":java_grpc_library.bzl", "java_grpc_library")
