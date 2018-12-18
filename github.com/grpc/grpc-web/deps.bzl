load("//:deps.bzl", 
    "com_github_grpc_grpc_web",
    "com_google_protobuf",
)

load("//closure:deps.bzl", 
    "io_bazel_rules_closure",
    "closure_proto_compile",
)

def closure_grpc_compile(**kwargs):
    com_google_protobuf(**kwargs)
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

