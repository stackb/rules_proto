# Aggregate all `dart` rules to one loadable file
load(":dart_proto_compile.bzl", _dart_proto_compile="dart_proto_compile")
load(":dart_grpc_compile.bzl", _dart_grpc_compile="dart_grpc_compile")
load(":dart_proto_library.bzl", _dart_proto_library="dart_proto_library")
load(":dart_grpc_library.bzl", _dart_grpc_library="dart_grpc_library")

dart_proto_compile = _dart_proto_compile
dart_grpc_compile = _dart_grpc_compile
dart_proto_library = _dart_proto_library
dart_grpc_library = _dart_grpc_library
