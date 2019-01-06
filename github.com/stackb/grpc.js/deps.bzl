load(
    "//:deps.bzl",
    "com_github_stackb_grpc_js",
    "com_google_protobuf",
    "io_bazel_rules_go",
)
load(
    "//closure:deps.bzl",
    "closure_proto_compile",
    "io_bazel_rules_closure",
)

def closure_grpc_compile(**kwargs):
    com_google_protobuf(**kwargs)
    io_bazel_rules_closure(**kwargs)
    io_bazel_rules_go(**kwargs)
    com_github_stackb_grpc_js(**kwargs)

def closure_grpc_library(**kwargs):
    closure_proto_compile(**kwargs)
    closure_grpc_compile(**kwargs)
