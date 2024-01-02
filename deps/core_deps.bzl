"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def core_deps():
    """core dependency macro
    """
    io_bazel_rules_go()  # via bazel_gazelle
    bazel_gazelle()  # via <TOP>
    rules_proto()  # via <TOP>

def io_bazel_rules_go():
    _maybe(
        http_archive,
        name = "io_bazel_rules_go",
        sha256 = "9bce3748e7ee751915d68a92021c48a03645f0028b7fbe33eb3560f456233934",
        strip_prefix = "rules_go-0.44.1",
        urls = [
            "https://github.com/bazelbuild/rules_go/archive/v0.44.1.tar.gz",
        ],
    )

def bazel_gazelle():
    _maybe(
        http_archive,
        name = "bazel_gazelle",
        sha256 = "a0ee1d304f7caa46680ba06bdef0e5d9ec8815f6e01ec29398efd13256598c3f",
        strip_prefix = "bazel-gazelle-0.35.0",
        urls = [
            "https://github.com/bazelbuild/bazel-gazelle/archive/v0.35.0.tar.gz",
        ],
        patches = [
            "@build_stack_rules_proto//third_party:bazel-gazelle-revert-1152.patch",
        ],
        patch_args = [
            "-p1",
        ],
    )

def rules_proto():
    _maybe(
        http_archive,
        name = "rules_proto",
        sha256 = "f5ae0e582238fcd4ea3d0146a3f5f3db9517f8fe24491eab3c105ace53aad1bb",
        strip_prefix = "rules_proto-f9b0b880d1e10e18daeeb168cef9d0f8316fdcb5",
        urls = [
            "https://github.com/bazelbuild/rules_proto/archive/f9b0b880d1e10e18daeeb168cef9d0f8316fdcb5.tar.gz",
        ],
    )
