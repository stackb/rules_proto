load("//:deps.bzl", 
    "com_github_grpc_grpc",
    "com_google_protobuf",
    "build_bazel_rules_swift",
)

def swift_proto_compile(**kwargs):
    com_google_protobuf(**kwargs)
    build_bazel_rules_swift()

def swift_grpc_compile(**kwargs):
    swift_proto_compile(**kwargs)

def swift_proto_library(**kwargs):
    swift_proto_compile(**kwargs)

def swift_grpc_library(**kwargs):
    swift_grpc_compile(**kwargs)
    swift_proto_library(**kwargs)