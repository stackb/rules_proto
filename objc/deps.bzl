load(
    "//:deps.bzl",
    "com_github_grpc_grpc",
)
load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def objc_proto_compile(**kwargs):
    protobuf(**kwargs)

def objc_grpc_compile(**kwargs):
    objc_proto_compile(**kwargs)
    com_github_grpc_grpc(**kwargs)
