load(
    "//:deps.bzl",
    "com_github_grpc_grpc",
    "io_bazel_rules_python",
    "six",
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
    io_bazel_rules_python(**kwargs)
    six(**kwargs)

def python_proto_library(**kwargs):
    python_proto_compile(**kwargs)
    io_bazel_rules_python(**kwargs)

def python_grpc_library(**kwargs):
    python_grpc_compile(**kwargs)
    python_proto_library(**kwargs)
