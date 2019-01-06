load(
    "//:deps.bzl",
    "com_google_guava_guava",
    "com_google_protobuf",
    "io_grpc_grpc_java",
    "javax_annotation_javax_annotation_api",
)

def java_proto_compile(**kwargs):
    com_google_protobuf(**kwargs)

def java_grpc_compile(**kwargs):
    java_proto_compile(**kwargs)
    io_grpc_grpc_java(**kwargs)

def java_proto_library(**kwargs):
    java_proto_compile(**kwargs)
    com_google_guava_guava(**kwargs)
    javax_annotation_javax_annotation_api(**kwargs)

def java_grpc_library(**kwargs):
    java_grpc_compile(**kwargs)
    java_proto_library(**kwargs)
