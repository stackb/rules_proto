load(
    "//:deps.bzl",
    "com_github_grpc_grpc",
    "com_github_yugui_rules_ruby",
)

load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def ruby_proto_compile(**kwargs):
    protobuf(**kwargs)

def ruby_grpc_compile(**kwargs):
    ruby_proto_compile(**kwargs)
    com_github_grpc_grpc(**kwargs)

def ruby_proto_library(**kwargs):
    ruby_proto_compile(**kwargs)
    com_github_yugui_rules_ruby(**kwargs)

def ruby_grpc_library(**kwargs):
    ruby_grpc_compile(**kwargs)
    ruby_proto_library(**kwargs)
