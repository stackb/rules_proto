load("//:deps.bzl",
    "com_github_grpc_grpc",
    "com_google_protobuf",
    "io_bazel_rules_dotnet",
)

def csharp_proto_compile(**kwargs):
    com_google_protobuf(**kwargs)

def csharp_grpc_compile(**kwargs):
    csharp_proto_compile(**kwargs)
    com_github_grpc_grpc(**kwargs)

def csharp_proto_library(**kwargs):
    csharp_proto_compile(**kwargs)
    io_bazel_rules_dotnet(**kwargs)

def csharp_grpc_library(**kwargs):
    csharp_grpc_compile(**kwargs)
    csharp_proto_library(**kwargs)
