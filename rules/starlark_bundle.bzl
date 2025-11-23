"""starlark_bundle.bzl aggreagates the transitive sources of dependencies

Rule is needed because bzl_library strictly requires srcs.
"""

load("@bazel_skylib//:bzl_library.bzl", "StarlarkLibraryInfo")

def _starlark_bundle_impl(ctx):
    deps = [d[StarlarkLibraryInfo] for d in ctx.attr.deps]
    transitive_srcs = depset(transitive = [d.transitive_srcs for d in deps])

    return [
        DefaultInfo(
            files = transitive_srcs,
        ),
        StarlarkLibraryInfo(
            srcs = [],
            transitive_srcs = transitive_srcs,
        ),
    ]

starlark_bundle = rule(
    implementation = _starlark_bundle_impl,
    attrs = {
        "deps": attr.label_list(
            doc = "list of starlark_library or bzl library rules",
            providers = [StarlarkLibraryInfo],
        ),
    },
    doc = "",
    provides = [DefaultInfo, StarlarkLibraryInfo],
)
