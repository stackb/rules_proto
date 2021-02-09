#
# Toolchains
#
def protoc_toolchain():
    # check_bazel_minimum_version(MINIMUM_BAZEL_VERSION)
    native.register_toolchains(str(Label("@build_stack_rules_proto//toolchains:protoc_toolchain")))
