load(
    "//:deps.bzl",
    "build_bazel_rules_android",
    "com_google_guava_guava_android",
    "com_google_protobuf",
    "io_grpc_grpc_java",
)
load(
    "//protobuf:deps.bzl",
    "protobuf_deps",
)

def android_deps(**kwargs):
    protobuf_deps(**kwargs)
    io_grpc_grpc_java(**kwargs)
    build_bazel_rules_android(**kwargs)
    com_google_guava_guava_android(**kwargs)

def android_proto_compile(**kwargs): # Kept for backwards compatibility
    android_deps(**kwargs)

def android_grpc_compile(**kwargs): # Kept for backwards compatibility
    android_deps(**kwargs)

def android_proto_library(**kwargs): # Kept for backwards compatibility
    android_deps(**kwargs)

def android_grpc_library(**kwargs): # Kept for backwards compatibility
    android_deps(**kwargs)
