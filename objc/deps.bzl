load(
    "//:deps.bzl",
    "com_github_grpc_grpc",
)
load(
    "//protobuf:deps.bzl",
    "protobuf_deps",
)

def objc_deps(**kwargs):
    protobuf_deps(**kwargs)
    com_github_grpc_grpc(**kwargs)

def objc_proto_compile(**kwargs): # Kept for backwards compatibility
    objc_deps(**kwargs)

def objc_proto_library(**kwargs): # Kept for backwards compatibility
    objc_deps(**kwargs)

def objc_grpc_compile(**kwargs): # Kept for backwards compatibility
    objc_deps(**kwargs)

def objc_grpc_library(**kwargs): # Kept for backwards compatibility
    objc_deps(**kwargs)
