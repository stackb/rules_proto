load(
    "//:deps.bzl",
    "com_github_grpc_grpc",
    "com_apt_itude_rules_pip",
    "six",
)
load(
    "//protobuf:deps.bzl",
    "protobuf_deps",
)

def python_deps(**kwargs):
    protobuf_deps(**kwargs)
    six(**kwargs)
    com_github_grpc_grpc(**kwargs)
    com_apt_itude_rules_pip(**kwargs)

def python_proto_compile(**kwargs): # Kept for backwards compatibility
    python_deps(**kwargs)

def python_grpc_compile(**kwargs): # Kept for backwards compatibility
    python_deps(**kwargs)

def python_proto_library(**kwargs): # Kept for backwards compatibility
    python_deps(**kwargs)

def python_grpc_library(**kwargs): # Kept for backwards compatibility
    python_deps(**kwargs)
