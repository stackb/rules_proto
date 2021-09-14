local_repository(
    name = "build_stack_rules_proto",
    path = "../build_stack_rules_proto",
)

register_toolchains("@build_stack_rules_proto//toolchain:prebuilt")

# == Externals ==

load("@build_stack_rules_proto//deps:core_deps.bzl", "core_deps")

core_deps()

# == Go ==

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.16.2")

# == Gazelle ==

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

# == Protobuf ==

load("@build_stack_rules_proto//deps:prebuilt_deps.bzl", "prebuilt_deps")

prebuilt_deps()
