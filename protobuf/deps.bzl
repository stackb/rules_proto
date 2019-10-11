load(
    "//:deps.bzl",
    "bazel_skylib",
    "com_github_madler_zlib",
    "com_google_protobuf",
    "external_protobuf_clib",
    "external_protobuf_headers",
    "rules_cc",
    "rules_java",
    "rules_proto",
    "rules_python",
)

def protobuf(**kwargs):
    com_github_madler_zlib(**kwargs)
    bazel_skylib(**kwargs)
    com_google_protobuf(**kwargs)
    external_protobuf_clib(**kwargs)
    external_protobuf_headers(**kwargs)
    # https://github.com/bazelbuild/rules_proto/blob/master/proto/private/dependencies.bzl
    rules_cc(**kwargs);
    rules_java(**kwargs);
    rules_proto(**kwargs);
    rules_python(**kwargs);
