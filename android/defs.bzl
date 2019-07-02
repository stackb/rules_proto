# Aggregate all `android` rules to one loadable file
load(":android_proto_compile.bzl", "android_proto_compile")
load(":android_grpc_compile.bzl", "android_grpc_compile")
load(":android_proto_library.bzl", "android_proto_library")
load(":android_grpc_library.bzl", "android_grpc_library")
