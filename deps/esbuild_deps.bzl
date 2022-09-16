"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_gazelle//:deps.bzl", "go_repository")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def esbuild_deps():
    com_github_evanw_esbuild()  # via <TOP>

def com_github_evanw_esbuild():
    _maybe(
        go_repository,
        name = "com_github_evanw_esbuild",
        sum = "h1:tWgtDpFR/VKWbWnSxigewZrxORvulw1Z6+KGOiJdNzA=",
        version = "v0.14.38",
        importpath = "github.com/evanw/esbuild",
        build_file_proto_mode = "disable_global",
    )
