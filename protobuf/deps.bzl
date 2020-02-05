load(
    "//:deps.bzl",
    "bazel_skylib",
    "zlib",
    "com_google_protobuf",
    "external_protobuf_clib",
    "external_protobuf_headers",
    "io_bazel_rules_python",
)

def protobuf(**kwargs):
    zlib(**kwargs)
    bazel_skylib(**kwargs)
    com_google_protobuf(**kwargs)
    external_protobuf_clib(**kwargs)
    external_protobuf_headers(**kwargs)
    io_bazel_rules_python(**kwargs)