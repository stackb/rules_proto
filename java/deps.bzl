load(
    "//:deps.bzl",
    "com_google_guava_guava",
    "com_google_protobuf",
    "io_grpc_grpc_java",
)

load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def java_proto_compile(**kwargs):
    protobuf(**kwargs)

def java_grpc_compile(**kwargs):
    java_proto_compile(**kwargs)
    io_grpc_grpc_java(**kwargs)

def java_proto_library(**kwargs):
    java_proto_compile(**kwargs)
    com_google_guava_guava(**kwargs)

def java_grpc_library(**kwargs):
    java_grpc_compile(**kwargs)
    java_proto_library(**kwargs)
