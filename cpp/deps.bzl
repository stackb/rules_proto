load(
    "//:deps.bzl",
    "com_github_grpc_grpc",
    "external_protobuf_clib",
    "external_madler_zlib",
    "bazel_version"
)
load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def cpp_proto_compile(**kwargs):
    protobuf(**kwargs)

def cpp_grpc_compile(**kwargs):
    bazel_version(**kwargs)
    cpp_proto_compile(**kwargs)
    external_madler_zlib(**kwargs)
    com_github_grpc_grpc(**kwargs)

def cpp_proto_library(**kwargs):
    cpp_proto_compile(**kwargs)
    external_protobuf_clib(**kwargs)

def cpp_grpc_library(**kwargs):
    cpp_grpc_compile(**kwargs)
    cpp_proto_library(**kwargs)
