load("//:deps.bzl", 
    "com_github_grpc_grpc",
    "com_google_protobuf", 
    "io_bazel_rules_dotnet",
)

load("//csharp/nuget:nuget.bzl", 
    "nuget_protobuf_packages", 
)

load("//csharp/nuget:nuget.bzl", 
    "nuget_grpc_packages", 
)

def csharp_proto_compile(**kwargs):
    com_google_protobuf(**kwargs)

def csharp_grpc_compile(**kwargs):
    csharp_proto_compile(**kwargs)
    com_github_grpc_grpc(**kwargs)

def csharp_proto_library(**kwargs):
    csharp_proto_compile(**kwargs)
    io_bazel_rules_dotnet(**kwargs)
    nuget_protobuf_packages(**kwargs)

def csharp_grpc_library(**kwargs):
    csharp_grpc_compile(**kwargs)
    csharp_proto_library(**kwargs)
    nuget_grpc_packages(**kwargs)
