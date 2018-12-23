load(
    "//:deps.bzl",
    "bazel_skylib",
    "build_bazel_rules_swift",
    "com_github_apple_swift_swift_protobuf",
    "com_github_grpc_grpc",
    "com_google_protobuf",
    "io_bazel_rules_go",
)

def swift_proto_compile(**kwargs):
    com_google_protobuf(**kwargs)
    io_bazel_rules_go(**kwargs)
    build_bazel_rules_swift(**kwargs)
    com_github_apple_swift_swift_protobuf(**kwargs)
    bazel_skylib(**kwargs)

def swift_grpc_compile(**kwargs):
    swift_proto_compile(**kwargs)

def swift_proto_library(**kwargs):
    swift_proto_compile(**kwargs)

def swift_grpc_library(**kwargs):
    swift_grpc_compile(**kwargs)
    swift_proto_library(**kwargs)
