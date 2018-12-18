load("//:deps.bzl", 
    "com_github_stackb_grpc_js",
    "com_google_protobuf",
    "io_bazel_rules_go",
)

load("//closure:deps.bzl", 
    "io_bazel_rules_closure",
    "closure_proto_compile",
)

def closure_grpc_compile(**kwargs):
    # protoc
    com_google_protobuf(**kwargs)
    # Need rules_go to build the plugin
    io_bazel_rules_go(**kwargs)
    # Need the plugin itself
    com_github_stackb_grpc_js(**kwargs)


def closure_grpc_library(**kwargs):
    closure_proto_compile(**kwargs)
    closure_grpc_compile(**kwargs)    
    io_bazel_rules_closure(**kwargs)
