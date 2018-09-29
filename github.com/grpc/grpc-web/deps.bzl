load("//:deps.bzl", 
    "com_github_grpc_grpc_web",
    "com_google_protobuf",
    "io_bazel_rules_closure",
)

load("//closure:deps.bzl", 
    "closure_proto_compile",
)

def web_grpc_compile(**kwargs):
    com_google_protobuf(**kwargs)
    com_github_grpc_grpc_web(**kwargs)


def web_grpc_library(**kwargs):
    closure_proto_compile(**kwargs)
    web_grpc_compile(**kwargs)    
    io_bazel_rules_closure(**kwargs)
