load(
    "//:deps.bzl",
    "io_bazel_rules_closure"
)
load(
    "//protobuf:deps.bzl",
    "protobuf_deps",
)

def closure_deps(**kwargs):
    protobuf_deps(**kwargs)
    io_bazel_rules_closure(**kwargs)

def closure_proto_compile(**kwargs): # Kept for backwards compatibility
    closure_deps(**kwargs)

def closure_proto_library(**kwargs): # Kept for backwards compatibility
    closure_deps(**kwargs)
