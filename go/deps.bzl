load("//:deps.bzl",
    "io_bazel_rules_go",
    "com_google_protobuf",
)

def go_proto_compile(**kwargs):
    com_google_protobuf(**kwargs)
    io_bazel_rules_go(**kwargs)

def go_grpc_compile(**kwargs):
    go_proto_compile(**kwargs)

def go_proto_library(**kwargs):
    go_proto_compile(**kwargs)

def go_grpc_library(**kwargs):
    go_grpc_compile(**kwargs)
    go_proto_library(**kwargs)
