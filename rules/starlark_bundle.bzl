"""starlark_bundle.bzl aggreagates the transitive sources of dependencies

Rule is needed because bzl_library strictly requires srcs.
"""

load("@bazel_skylib//:bzl_library.bzl", "StarlarkLibraryInfo")
load(":starlark_library.bzl", "StarlarkLibraryFileInfo")

def _starlark_bundle_impl(ctx):
    deps = [d[StarlarkLibraryFileInfo] for d in ctx.attr.deps]
    transitive_srcs = depset(transitive = [d.transitive_srcs for d in deps])
    transitive_docs = depset(transitive = [d.transitive_docs for d in deps])

    return [
        DefaultInfo(
            files = transitive_docs,
        ),
        StarlarkLibraryInfo(
            srcs = [],
            transitive_srcs = transitive_srcs,
        ),
        StarlarkLibraryFileInfo(
            src = None,
            transitive_srcs = transitive_srcs,
            transitive_docs = transitive_docs,
            broken = False,
        ),
    ]

_starlark_bundle = rule(
    implementation = _starlark_bundle_impl,
    attrs = {
        "deps": attr.label_list(
            doc = "list of starlark_library or bzl library rules",
            providers = [StarlarkLibraryFileInfo],
        ),
    },
    doc = "",
    provides = [DefaultInfo, StarlarkLibraryFileInfo, StarlarkLibraryInfo],
)

def starlark_bundle(name, deps, **kwargs):
    _starlark_bundle(
        name = name,
        deps = deps,
        **kwargs
    )
