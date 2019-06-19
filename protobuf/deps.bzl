load(
    "//:deps.bzl",
    "bazel_skylib",
    "com_google_protobuf",
    "external_zlib",
)

def protobuf(**kwargs):
    bazel_skylib(**kwargs)
    com_google_protobuf(**kwargs)
    external_zlib(**kwargs)
    native.register_toolchains(str(Label("//protobuf:protoc_toolchain")))
