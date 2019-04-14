load(
    "//:deps.bzl",
    "com_github_stackb_grpc_js",
    "io_bazel_rules_go",
)
load(
    "//closure:deps.bzl",
    "closure_proto_compile",
    "io_bazel_rules_closure",
)

load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def closure_grpc_compile(**kwargs):
    protobuf(**kwargs)
    io_bazel_rules_go(**kwargs)
    io_bazel_rules_closure(**kwargs)
    com_github_stackb_grpc_js(**kwargs)

def closure_grpc_library(**kwargs):
    closure_proto_compile(**kwargs)
    closure_grpc_compile(**kwargs)
    io_bazel_rules_closure(**kwargs)
