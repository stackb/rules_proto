# Aggregate all `objc` rules to one loadable file
load(":objc_proto_compile.bzl", _objc_proto_compile="objc_proto_compile")
load(":objc_grpc_compile.bzl", _objc_grpc_compile="objc_grpc_compile")
load(":objc_proto_library.bzl", _objc_proto_library="objc_proto_library")

objc_proto_compile = _objc_proto_compile
objc_grpc_compile = _objc_grpc_compile
objc_proto_library = _objc_proto_library
