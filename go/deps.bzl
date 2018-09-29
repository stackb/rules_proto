load("//:deps.bzl", 
    "io_bazel_rules_go",
)

def go_proto_compile(**kwargs):
    io_bazel_rules_go(**kwargs)

def go_grpc_compile(**kwargs):
    go_proto_compile(**kwargs)

def go_proto_library(**kwargs):
    go_proto_compile(**kwargs)

def go_grpc_library(**kwargs):
    go_grpc_compile(**kwargs)
    go_proto_library(**kwargs)
