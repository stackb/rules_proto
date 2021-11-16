"""
GENERATED FILE - DO NOT EDIT (created via @build_stack_rules_proto//cmd/depsgen)
"""

load("@bazel_gazelle//:deps.bzl", "go_repository")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def go_core_deps():
    com_github_golang_protobuf()  # via <TOP>
    org_golang_google_grpc()  # via <TOP>
    org_golang_google_grpc_cmd_protoc_gen_go_grpc()  # via <TOP>

def com_github_golang_protobuf():
    _maybe(
        go_repository,
        name = "com_github_golang_protobuf",
        sum = "h1:JjCZWpVbqXDqFVmTfYWEVTMIYrL/NPdPSCHPJ0T/raM=",
        version = "v1.4.3",
        importpath = "github.com/golang/protobuf",
        build_file_proto_mode = "disable_global",
    )

def org_golang_google_grpc():
    _maybe(
        go_repository,
        name = "org_golang_google_grpc",
        sum = "h1:TwIQcH3es+MojMVojxxfQ3l3OF2KzlRxML2xZq0kRo8=",
        version = "v1.35.0",
        importpath = "google.golang.org/grpc",
        build_file_proto_mode = "disable",
    )

def org_golang_google_grpc_cmd_protoc_gen_go_grpc():
    _maybe(
        go_repository,
        name = "org_golang_google_grpc_cmd_protoc_gen_go_grpc",
        sum = "h1:M1YKkFIboKNieVO5DLUEVzQfGwJD30Nv2jfUgzb5UcE=",
        version = "v1.1.0",
        importpath = "google.golang.org/grpc/cmd/protoc-gen-go-grpc",
        build_file_proto_mode = "disable_global",
    )
