load(
    "//:deps.bzl",
    "com_github_grpc_grpc",
)
load(
    "//protobuf:deps.bzl",
    "protobuf_deps",
)

def cpp_deps(**kwargs):
    protobuf_deps(**kwargs)
    com_github_grpc_grpc(**kwargs)

def cpp_proto_compile(**kwargs): # Kept for backwards compatibility
    cpp_deps(**kwargs)

def cpp_grpc_compile(**kwargs): # Kept for backwards compatibility
    cpp_deps(**kwargs)

def cpp_proto_library(**kwargs): # Kept for backwards compatibility
    cpp_deps(**kwargs)

def cpp_grpc_library(**kwargs): # Kept for backwards compatibility
    cpp_deps(**kwargs)
