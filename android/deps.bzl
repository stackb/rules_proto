load(
    "//:deps.bzl",
    "build_bazel_rules_android",
    "com_google_guava_guava_android",
    "com_google_protobuf",
    "com_google_protobuf_lite",
    "gmaven_rules",
    "io_grpc_grpc_java",
)
load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def android_proto_compile(**kwargs):
    protobuf(**kwargs)
    com_google_protobuf_lite(**kwargs)

def android_grpc_compile(**kwargs):
    android_proto_compile(**kwargs)
    io_grpc_grpc_java(**kwargs)

def android_proto_library(**kwargs):
    android_proto_compile(**kwargs)
    build_bazel_rules_android(**kwargs)
    gmaven_rules(**kwargs)
    com_google_guava_guava_android(**kwargs)

def android_grpc_library(**kwargs):
    android_grpc_compile(**kwargs)
    android_proto_library(**kwargs)
