load(
    "//:deps.bzl",
    "com_github_grpc_grpc",
)

load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def python_proto_compile(**kwargs):
    protobuf(**kwargs)

def python_grpc_compile(**kwargs):
    python_proto_compile(**kwargs)
    com_github_grpc_grpc(**kwargs)

def python_proto_library(**kwargs):
    python_proto_compile(**kwargs)

def python_grpc_library(**kwargs):
    python_grpc_compile(**kwargs)
    python_proto_library(**kwargs)
