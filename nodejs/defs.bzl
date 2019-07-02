# Aggregate all `nodejs` rules to one loadable file
load(":nodejs_proto_compile.bzl", _nodejs_proto_compile="nodejs_proto_compile")
load(":nodejs_grpc_compile.bzl", _nodejs_grpc_compile="nodejs_grpc_compile")

nodejs_proto_compile = _nodejs_proto_compile
nodejs_grpc_compile = _nodejs_grpc_compile
