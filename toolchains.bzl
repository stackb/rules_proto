
#
# Toolchains
#
def protoc_toolchains():
    # check_bazel_minimum_version(MINIMUM_BAZEL_VERSION)
    native.register_toolchains(str(Label("//protobuf:protoc_toolchain")))
