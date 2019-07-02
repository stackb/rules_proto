# Aggregate all `python` rules to one loadable file
load(":python_proto_compile.bzl", _python_proto_compile="python_proto_compile")
load(":python_grpc_compile.bzl", _python_grpc_compile="python_grpc_compile")
load(":python_proto_library.bzl", _python_proto_library="python_proto_library")
load(":python_grpc_library.bzl", _python_grpc_library="python_grpc_library")

python_proto_compile = _python_proto_compile
python_grpc_compile = _python_grpc_compile
python_proto_library = _python_proto_library
python_grpc_library = _python_grpc_library
