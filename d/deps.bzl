load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load(
    "//:deps.bzl",
    "io_bazel_rules_d",
)
load(
    "//protobuf:deps.bzl",
    "protobuf",
)

def com_github_dcarp_protobuf_d():
    http_archive(
        name = "com_github_dcarp_protobuf_d",
        urls = ["https://github.com/dcarp/protobuf-d/archive/v0.5.0.tar.gz"],
        strip_prefix = "protobuf-d-0.5.0",
        sha256 = "67a037dc29242f0d2f099746da67f40afff27c07f9ab48dda53d5847620db421",
        build_file = Label("//d:com_github_dcarp_protobuf_d.BUILD.bazel"),
    )

def d_proto_compile(**kwargs):
    protobuf(**kwargs)
    com_github_dcarp_protobuf_d()
    io_bazel_rules_d(**kwargs)

def d_grpc_compile(**kwargs):
    d_proto_compile(**kwargs)

def d_proto_library(**kwargs):
    d_proto_compile(**kwargs)

def d_grpc_library(**kwargs):
    d_grpc_compile(**kwargs)
    d_proto_library(**kwargs)
