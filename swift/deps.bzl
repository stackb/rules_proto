load(
    "//:deps.bzl",
    "build_bazel_rules_swift",
    "io_bazel_rules_go",
    "com_github_apple_swift_swift_protobuf",
    "com_github_grpc_grpc",
)

load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def swift_proto_compile(**kwargs):
    protobuf(**kwargs)
    # io_bazel_rules_go(**kwargs)
    build_bazel_rules_swift(**kwargs)

def swift_grpc_compile(**kwargs):
    protobuf(**kwargs)
    # io_bazel_rules_go(**kwargs)
    build_bazel_rules_swift(**kwargs)

def swift_proto_library(**kwargs):
    build_bazel_rules_swift(**kwargs)

def swift_grpc_library(**kwargs):
    build_bazel_rules_swift(**kwargs)
