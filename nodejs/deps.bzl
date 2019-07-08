load(
    "//:deps.bzl",
    "com_github_grpc_grpc",
    "build_bazel_rules_nodejs",
)
load(
    "//protobuf:deps.bzl",
    "protobuf_deps",
)

def nodejs_deps(**kwargs):
    protobuf_deps(**kwargs)
    com_github_grpc_grpc(**kwargs)
    build_bazel_rules_nodejs(**kwargs)

def nodejs_proto_compile(**kwargs): # Kept for backwards compatibility
    nodejs_deps(**kwargs)

def nodejs_grpc_compile(**kwargs): # Kept for backwards compatibility
    nodejs_deps(**kwargs)

def nodejs_proto_library(**kwargs): # Kept for backwards compatibility
    nodejs_deps(**kwargs)

def nodejs_grpc_library(**kwargs): # Kept for backwards compatibility
    nodejs_deps(**kwargs)
