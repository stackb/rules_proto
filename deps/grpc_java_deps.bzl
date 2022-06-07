"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def grpc_java_deps():
    rules_jvm_external()  # via io_grpc_grpc_java
    io_grpc_grpc_java()  # via <TOP>

def rules_jvm_external():
    _maybe(
        http_archive,
        name = "rules_jvm_external",
        sha256 = "31701ad93dbfe544d597dbe62c9a1fdd76d81d8a9150c2bf1ecf928ecdf97169",
        strip_prefix = "rules_jvm_external-4.0",
        urls = [
            "https://github.com/bazelbuild/rules_jvm_external/archive/4.0.zip",
        ],
    )

def io_grpc_grpc_java():
    _maybe(
        http_archive,
        name = "io_grpc_grpc_java",
        sha256 = "82b3cf09f98a5932e1b55175aaec91b2a3f424eec811e47b2a3be533044d9afb",
        strip_prefix = "grpc-java-7f7821c616598ce4e33d2045c5641b2348728cb8",
        urls = [
            "https://github.com/grpc/grpc-java/archive/7f7821c616598ce4e33d2045c5641b2348728cb8.tar.gz",
        ],
    )
