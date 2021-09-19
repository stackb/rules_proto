"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_gazelle//:deps.bzl", "go_repository")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def go_core_deps():
    org_golang_google_grpc()  # via <TOP>

def org_golang_google_grpc():
    _maybe(
        go_repository,
        name = "org_golang_google_grpc",
        sum = "h1:TwIQcH3es+MojMVojxxfQ3l3OF2KzlRxML2xZq0kRo8=",
        version = "v1.35.0",
        importpath = "google.golang.org/grpc",
        build_file_proto_mode = "disable",
    )
