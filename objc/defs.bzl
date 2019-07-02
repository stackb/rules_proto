# Aggregate all `objc` rules to one loadable file
load(":objc_proto_compile.bzl", "objc_proto_compile")
load(":objc_grpc_compile.bzl", "objc_grpc_compile")
load(":objc_proto_library.bzl", "objc_proto_library")
