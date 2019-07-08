load(
    "//:deps.bzl",
    "com_github_grpc_grpc",
    "com_github_yugui_rules_ruby",
)
load(
    "//protobuf:deps.bzl",
    "protobuf_deps",
)

def ruby_deps(**kwargs):
    protobuf_deps(**kwargs)
    com_github_grpc_grpc(**kwargs)
    com_github_yugui_rules_ruby(**kwargs)

def ruby_proto_compile(**kwargs): # Kept for backwards compatibility
    ruby_deps(**kwargs)

def ruby_grpc_compile(**kwargs): # Kept for backwards compatibility
    ruby_deps(**kwargs)

def ruby_proto_library(**kwargs): # Kept for backwards compatibility
    ruby_deps(**kwargs)

def ruby_grpc_library(**kwargs): # Kept for backwards compatibility
    ruby_deps(**kwargs)
