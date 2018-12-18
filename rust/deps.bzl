load("//:deps.bzl",
    "com_google_protobuf",
    "io_bazel_rules_rust",
)

def rust_proto_compile(**kwargs):
    com_google_protobuf(**kwargs)
    io_bazel_rules_rust(**kwargs)

def rust_grpc_compile(**kwargs):
    rust_proto_compile(**kwargs)

def rust_proto_library(**kwargs):
    rust_proto_compile(**kwargs)
    io_bazel_rules_rust(**kwargs)

def rust_grpc_library(**kwargs):
    rust_grpc_compile(**kwargs)
    rust_proto_library(**kwargs)
