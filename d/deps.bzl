load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load(
    "//:deps.bzl",
    "io_bazel_rules_d",
    "com_github_dcarp_protobuf_d",
)
load(
    "//protobuf:deps.bzl",
    "protobuf_deps",
)

def d_deps(**kwargs):
    protobuf_deps(**kwargs)
    com_github_dcarp_protobuf_d(**kwargs)
    io_bazel_rules_d(**kwargs)

def d_proto_compile(**kwargs): # Kept for backwards compatibility
    d_deps(**kwargs)

def d_grpc_compile(**kwargs): # Kept for backwards compatibility
    d_deps(**kwargs)

def d_proto_library(**kwargs): # Kept for backwards compatibility
    d_deps(**kwargs)

def d_grpc_library(**kwargs): # Kept for backwards compatibility
    d_deps(**kwargs)
