"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def grpc_java_deps():
    """grpc_java dependency macro
    """
    rules_jvm_external()  # via io_grpc_grpc_java
    io_grpc_grpc_java()  # via <TOP>

def rules_jvm_external():
    _maybe(
        http_archive,
        name = "rules_jvm_external",
        sha256 = "1ce86ffee65725300dc1f0017b7df89715c832de550137432dc1985d60a13155",
        strip_prefix = "rules_jvm_external-e6c1ff21e002bf97a7b1c07d63edd508a8dc9659",
        urls = [
            "https://github.com/bazelbuild/rules_jvm_external/archive/e6c1ff21e002bf97a7b1c07d63edd508a8dc9659.tar.gz",
        ],
    )

def io_grpc_grpc_java():
    _maybe(
        http_archive,
        name = "io_grpc_grpc_java",
        sha256 = "4a021ea9ebb96f5841a135c168209cf2413587a0f8ce71a2bf37b3aad847b0d0",
        strip_prefix = "grpc-java-1.57.1",
        urls = [
            "https://github.com/grpc/grpc-java/archive/v1.57.1.tar.gz",
        ],
    )
