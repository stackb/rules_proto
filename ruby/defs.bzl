# Aggregate all `ruby` rules to one loadable file
load(":ruby_proto_compile.bzl", "ruby_proto_compile")
load(":ruby_grpc_compile.bzl", "ruby_grpc_compile")
load(":ruby_proto_library.bzl", "ruby_proto_library")
load(":ruby_grpc_library.bzl", "ruby_grpc_library")
