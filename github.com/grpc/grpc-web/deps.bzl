load(
    "//:deps.bzl",
    "com_github_grpc_grpc_web",
)
load(
    "//closure:deps.bzl",
    "closure_deps",
)

def grpc_web_deps(**kwargs):
    closure_deps(**kwargs)
    com_github_grpc_grpc_web(**kwargs)

def closure_grpc_compile(**kwargs): # Kept for backwards compatibility
    grpc_web_deps(**kwargs)

def commonjs_grpc_compile(**kwargs): # Kept for backwards compatibility
    grpc_web_deps(**kwargs)

def commonjs_dts_grpc_compile(**kwargs): # Kept for backwards compatibility
    grpc_web_deps(**kwargs)

def ts_grpc_compile(**kwargs): # Kept for backwards compatibility
    grpc_web_deps(**kwargs)

def closure_grpc_library(**kwargs): # Kept for backwards compatibility
    grpc_web_deps(**kwargs)
