"""starlark_library.bzl is a thin wrapper over bzl_library."""

load("@bazel_skylib//:bzl_library.bzl", "bzl_library")

def starlark_library(name, **kwargs):
    bzl_library(
        name = name,
        **kwargs
    )
