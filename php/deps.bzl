load(
    "//:deps.bzl",
    "com_github_grpc_grpc",
)
load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def php_proto_compile(**kwargs):
    protobuf(**kwargs)

def php_grpc_compile(**kwargs):
    php_proto_compile(**kwargs)
    com_github_grpc_grpc(**kwargs)
