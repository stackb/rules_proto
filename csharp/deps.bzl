load(
    "//:deps.bzl",
    "com_github_grpc_grpc",
    "io_bazel_rules_dotnet",
)
load(
    "//protobuf:deps.bzl",
    "protobuf_deps",
)

def csharp_deps(**kwargs):
    protobuf_deps(**kwargs)
    com_github_grpc_grpc(**kwargs)
    io_bazel_rules_dotnet(**kwargs)

def csharp_proto_compile(**kwargs): # Kept for backwards compatibility
    csharp_deps(**kwargs)

def csharp_grpc_compile(**kwargs): # Kept for backwards compatibility
    csharp_deps(**kwargs)

def csharp_proto_library(**kwargs): # Kept for backwards compatibility
    csharp_deps(**kwargs)

def csharp_grpc_library(**kwargs): # Kept for backwards compatibility
    csharp_deps(**kwargs)
