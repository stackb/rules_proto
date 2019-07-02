# Aggregate all `go` rules to one loadable file
load(":go_proto_compile.bzl", "go_proto_compile")
load(":go_grpc_compile.bzl", "go_grpc_compile")
load(":go_proto_library.bzl", "go_proto_library")
load(":go_grpc_library.bzl", "go_grpc_library")
