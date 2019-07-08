load(
    "//:deps.bzl",
    "io_bazel_rules_go",
)
load(
    "//protobuf:deps.bzl",
    "protobuf_deps",
)

def go_deps(**kwargs):
    protobuf_deps(**kwargs)
    io_bazel_rules_go(**kwargs)

def go_proto_compile(**kwargs): # Kept for backwards compatibility
    go_deps(**kwargs)

def go_grpc_compile(**kwargs): # Kept for backwards compatibility
    go_deps(**kwargs)

def go_proto_library(**kwargs): # Kept for backwards compatibility
    go_deps(**kwargs)

def go_grpc_library(**kwargs): # Kept for backwards compatibility
    go_deps(**kwargs)
