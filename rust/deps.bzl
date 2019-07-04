load(
    "//:deps.bzl",
    "io_bazel_rules_rust",
)
load(
    "//protobuf:deps.bzl",
    "protobuf",
)
load(
    "@build_stack_rules_proto//rust/raze:crates.bzl",
    "raze_fetch_remote_crates"
)

def rust_proto_compile(**kwargs):
    protobuf(**kwargs)
    io_bazel_rules_rust(**kwargs)

def rust_grpc_compile(**kwargs):
    rust_proto_compile(**kwargs)

def rust_proto_library(**kwargs):
    rust_proto_compile(**kwargs)
    io_bazel_rules_rust(**kwargs)
    raze_fetch_remote_crates()

def rust_grpc_library(**kwargs):
    rust_grpc_compile(**kwargs)
    rust_proto_library(**kwargs)
