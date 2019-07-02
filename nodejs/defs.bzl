# Aggregate all `nodejs` rules to one loadable file
load(":nodejs_proto_compile.bzl", "nodejs_proto_compile")
load(":nodejs_grpc_compile.bzl", "nodejs_grpc_compile")
