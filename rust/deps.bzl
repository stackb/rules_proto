load(
    "//:deps.bzl",
    "io_bazel_rules_rust",
)
load(
    "//protobuf:deps.bzl",
    "protobuf_deps",
)
load(
    "//rust/raze:crates.bzl",
    "raze_fetch_remote_crates"
)

def rust_deps(**kwargs):
    protobuf_deps(**kwargs)
    io_bazel_rules_rust(**kwargs)
    raze_fetch_remote_crates()

def rust_proto_compile(**kwargs): # Kept for backwards compatibility
    rust_deps(**kwargs)

def rust_grpc_compile(**kwargs): # Kept for backwards compatibility
    rust_deps(**kwargs)

def rust_proto_library(**kwargs): # Kept for backwards compatibility
    rust_deps(**kwargs)

def rust_grpc_library(**kwargs): # Kept for backwards compatibility
    rust_deps(**kwargs)
