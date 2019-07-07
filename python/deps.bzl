load(
    "//:deps.bzl",
    "com_github_grpc_grpc",
    "com_apt_itude_rules_pip",
    "six",
)
load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def python_proto_compile(**kwargs):
    protobuf(**kwargs)
    six(**kwargs)

def python_grpc_compile(**kwargs):
    python_proto_compile(**kwargs)
    com_github_grpc_grpc(**kwargs)
    com_apt_itude_rules_pip(**kwargs)

def python_proto_library(**kwargs):
    python_proto_compile(**kwargs)
    com_apt_itude_rules_pip(**kwargs)

def python_grpc_library(**kwargs):
    python_grpc_compile(**kwargs)
    python_proto_library(**kwargs)
