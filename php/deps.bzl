load("//:deps.bzl", 
    "com_github_grpc_grpc",
    "com_google_protobuf", 
)

def php_proto_compile(**kwargs):
    com_google_protobuf(**kwargs)

def php_grpc_compile(**kwargs):
    php_proto_compile(**kwargs)
    com_github_grpc_grpc(**kwargs)
