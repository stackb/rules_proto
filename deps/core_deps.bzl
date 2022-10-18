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
        sha256 = "cc027f11f98aef8bc52c472ced0714994507a16ccd3a0820b2df2d6db695facd",
        strip_prefix = "rules_go-0.35.0",
        urls = [
            "https://github.com/bazelbuild/rules_go/archive/v0.35.0.tar.gz",
        ],
    )

def bazel_gazelle():
    _maybe(
        http_archive,
        name = "bazel_gazelle",
        sha256 = "fd86cd74ceaa68f0a6998c255d5f570a6cc20c345652da9d757352ac9e8c6cc9",
        strip_prefix = "bazel-gazelle-ba090ac841ca2a1211a21ce3dee38ea597b00099",
        urls = [
            "https://github.com/bazelbuild/bazel-gazelle/archive/ba090ac841ca2a1211a21ce3dee38ea597b00099.tar.gz",
        ],
        patches = [
            "@build_stack_rules_proto//third_party:bazel-gazelle-PR1274.patch",
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
        sha256 = "9fc210a34f0f9e7cc31598d109b5d069ef44911a82f507d5a88716db171615a8",
        strip_prefix = "rules_proto-f7a30f6f80006b591fa7c437fe5a951eb10bcbcf",
        urls = [
            "https://github.com/bazelbuild/rules_proto/archive/f7a30f6f80006b591fa7c437fe5a951eb10bcbcf.tar.gz",
        ],
    )
