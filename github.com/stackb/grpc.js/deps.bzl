load(
    "//:deps.bzl",
    "com_github_stackb_grpc_js",
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

def grpcjs_grpc_compile(**kwargs):
    protobuf(**kwargs)
    com_github_stackb_grpc_js(**kwargs)

def grpcjs_grpc_library(**kwargs):
    closure_proto_compile(**kwargs)
    grpcjs_grpc_compile(**kwargs)
    io_bazel_rules_closure(**kwargs)
