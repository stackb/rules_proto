load("//:deps.bzl", 
    "com_github_grpc_grpc",
    "com_google_protobuf", 
)

def cpp_proto_compile(**kwargs):
    com_google_protobuf(**kwargs)

def cpp_grpc_compile(**kwargs):
    cpp_proto_compile(**kwargs)
    com_github_grpc_grpc(**kwargs)

def cpp_proto_library(**kwargs):
    cpp_proto_compile(**kwargs)

def cpp_grpc_library(**kwargs):
    cpp_grpc_compile(**kwargs)
    cpp_proto_library(**kwargs)
