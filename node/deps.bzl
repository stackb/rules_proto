load(
    "//:deps.bzl",
    "com_github_grpc_grpc",
    "build_bazel_rules_nodejs",
)
load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def node_proto_compile(**kwargs):
    protobuf(**kwargs)

def node_grpc_compile(**kwargs):
    node_proto_compile(**kwargs)
    com_github_grpc_grpc(**kwargs)

def node_proto_library(**kwargs):
    node_proto_compile(**kwargs)
    build_bazel_rules_nodejs(**kwargs)

def node_grpc_library(**kwargs):
    node_grpc_compile(**kwargs)
    node_proto_library(**kwargs)
