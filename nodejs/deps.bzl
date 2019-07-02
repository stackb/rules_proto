load(
    "//:deps.bzl",
    "com_github_grpc_grpc",
    "build_bazel_rules_nodejs",
)
load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def nodejs_proto_compile(**kwargs):
    protobuf(**kwargs)

def nodejs_grpc_compile(**kwargs):
    nodejs_proto_compile(**kwargs)
    com_github_grpc_grpc(**kwargs)

def nodejs_proto_library(**kwargs):
    nodejs_proto_compile(**kwargs)
    build_bazel_rules_nodejs(**kwargs)

def nodejs_grpc_library(**kwargs):
    nodejs_grpc_compile(**kwargs)
    nodejs_proto_library(**kwargs)
