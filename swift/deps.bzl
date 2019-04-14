load(
    "//:deps.bzl",
    "build_bazel_rules_swift",
)

load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def swift_proto_compile(**kwargs):
    protobuf(**kwargs)
    build_bazel_rules_swift(**kwargs)

def swift_grpc_compile(**kwargs):
    protobuf(**kwargs)
    build_bazel_rules_swift(**kwargs)

def swift_proto_library(**kwargs):
    build_bazel_rules_swift(**kwargs)

def swift_grpc_library(**kwargs):
    build_bazel_rules_swift(**kwargs)
