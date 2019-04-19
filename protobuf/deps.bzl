load(
    "//:deps.bzl",
    "external_zlib",
    "external_protobuf_clib",
    "external_protobuf_headers",
    "bazel_skylib",
    "com_github_madler_zlib",
    "com_google_protobuf",
)

def protobuf(**kwargs):
    com_github_madler_zlib(**kwargs)
    external_zlib(**kwargs)
    bazel_skylib(**kwargs)
    com_google_protobuf(**kwargs)
    external_protobuf_clib(**kwargs)
    external_protobuf_headers(**kwargs)

