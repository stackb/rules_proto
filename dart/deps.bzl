
load("//:deps.bzl", 
    "com_google_protobuf",
    "dart_pub_deps_protoc_plugin",
    "dart_sdk",
    "io_bazel_rules_dart",
)

def dart_proto_compile(**kwargs):
    com_google_protobuf(**kwargs)
    io_bazel_rules_dart(**kwargs)
    dart_sdk(**kwargs)
    dart_pub_deps_protoc_plugin(**kwargs)

def dart_grpc_compile(**kwargs):
    dart_proto_compile(**kwargs)

def dart_proto_library(**kwargs):
    dart_proto_compile(**kwargs)

def dart_grpc_library(**kwargs):
    dart_grpc_compile(**kwargs)
    dart_proto_library(**kwargs)

