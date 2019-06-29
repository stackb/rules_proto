load(
    "//:deps.bzl",
    "com_github_grpc_grpc_web",
    "io_bazel_rules_closure",
)
load(
    "//closure:deps.bzl",
    "closure_proto_compile",
)

def closure_grpc_compile(**kwargs):
    closure_proto_compile(**kwargs)
    io_bazel_rules_closure(**kwargs)
    com_github_grpc_grpc_web(**kwargs)

def commonjs_grpc_compile(**kwargs):
    closure_grpc_compile(**kwargs)

def commonjs_dts_grpc_compile(**kwargs):
    closure_grpc_compile(**kwargs)

def ts_grpc_compile(**kwargs):
    closure_grpc_compile(**kwargs)

def closure_grpc_library(**kwargs):
    closure_grpc_compile(**kwargs)
