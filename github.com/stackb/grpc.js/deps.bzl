load(
    "//:deps.bzl",
    "com_github_stackb_grpc_js",
)
load(
    "//closure:deps.bzl",
    "closure_deps",
)

def grpcjs_deps(**kwargs):
    closure_deps(**kwargs)
    com_github_stackb_grpc_js(**kwargs)

def grpcjs_grpc_compile(**kwargs): # Kept for backwards compatibility
    grpcjs_deps(**kwargs)

def grpcjs_grpc_library(**kwargs): # Kept for backwards compatibility
    grpcjs_deps(**kwargs)
