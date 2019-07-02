# Aggregate all `gogo` rules to one loadable file
load(":gogo_proto_compile.bzl", _gogo_proto_compile="gogo_proto_compile")
load(":gogo_grpc_compile.bzl", _gogo_grpc_compile="gogo_grpc_compile")
load(":gogo_proto_library.bzl", _gogo_proto_library="gogo_proto_library")
load(":gogo_grpc_library.bzl", _gogo_grpc_library="gogo_grpc_library")
load(":gogofast_proto_compile.bzl", _gogofast_proto_compile="gogofast_proto_compile")
load(":gogofast_grpc_compile.bzl", _gogofast_grpc_compile="gogofast_grpc_compile")
load(":gogofast_proto_library.bzl", _gogofast_proto_library="gogofast_proto_library")
load(":gogofast_grpc_library.bzl", _gogofast_grpc_library="gogofast_grpc_library")
load(":gogofaster_proto_compile.bzl", _gogofaster_proto_compile="gogofaster_proto_compile")
load(":gogofaster_grpc_compile.bzl", _gogofaster_grpc_compile="gogofaster_grpc_compile")
load(":gogofaster_proto_library.bzl", _gogofaster_proto_library="gogofaster_proto_library")
load(":gogofaster_grpc_library.bzl", _gogofaster_grpc_library="gogofaster_grpc_library")

gogo_proto_compile = _gogo_proto_compile
gogo_grpc_compile = _gogo_grpc_compile
gogo_proto_library = _gogo_proto_library
gogo_grpc_library = _gogo_grpc_library
gogofast_proto_compile = _gogofast_proto_compile
gogofast_grpc_compile = _gogofast_grpc_compile
gogofast_proto_library = _gogofast_proto_library
gogofast_grpc_library = _gogofast_grpc_library
gogofaster_proto_compile = _gogofaster_proto_compile
gogofaster_grpc_compile = _gogofaster_grpc_compile
gogofaster_proto_library = _gogofaster_proto_library
gogofaster_grpc_library = _gogofaster_grpc_library
