load(
    "//:deps.bzl",
    "build_bazel_rules_swift",
)
load(
    "//protobuf:deps.bzl",
    "protobuf_deps",
)

def swift_deps(**kwargs):
    protobuf_deps(**kwargs)
    build_bazel_rules_swift(**kwargs)

def swift_proto_compile(**kwargs): # Kept for backwards compatibility
    swift_deps(**kwargs)

def swift_grpc_compile(**kwargs): # Kept for backwards compatibility
    swift_deps(**kwargs)

def swift_proto_library(**kwargs): # Kept for backwards compatibility
    swift_deps(**kwargs)

def swift_grpc_library(**kwargs): # Kept for backwards compatibility
    swift_deps(**kwargs)
