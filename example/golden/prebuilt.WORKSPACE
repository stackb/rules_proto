local_repository(
    name = "build_stack_rules_proto",
    path = "../build_stack_rules_proto",
)

register_toolchains("@build_stack_rules_proto//toolchain:prebuilt")

# == Externals ==

load("@build_stack_rules_proto//deps:core_deps.bzl", "core_deps")

core_deps()

# NOTE: these come before go_rules_dependencies and gazelle_dependencies in
# order to use newer grpc deps.

load("@build_stack_rules_proto//deps:go_core_deps.bzl", "go_core_deps")

go_core_deps()

# == Go ==

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.16.2")

# == Gazelle ==

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

load("@build_stack_rules_proto//:go_deps.bzl", "gazelle_protobuf_extension_go_deps")

gazelle_protobuf_extension_go_deps()

# == Protobuf ==

load("@build_stack_rules_proto//deps:prebuilt_protoc_deps.bzl", "prebuilt_protoc_deps")

prebuilt_protoc_deps()

load("@build_stack_rules_proto//deps:protobuf_core_deps.bzl", "protobuf_core_deps")

protobuf_core_deps()

# == gRPC Gateway ==

load("@build_stack_rules_proto//deps:grpc_gateway_deps.bzl", "grpc_gateway_deps")

grpc_gateway_deps()

load("@build_stack_rules_proto//deps:grpc_gateway_openapiv2_deps.bzl", "grpc_gateway_openapiv2_deps")

grpc_gateway_openapiv2_deps()
