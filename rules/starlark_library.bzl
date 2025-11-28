"""starlark_library.bzl is a thin wrapper over bzl_library."""

load("@bazel_skylib//:bzl_library.bzl", "StarlarkLibraryInfo")

StarlarkLibraryFileInfo = provider(
    "Information on contained Starlark rules.",
    fields = {
        "label": "The label of the target rule",
        "repo_name": "The original (non-canonical) repo (workspace!) name",
        "src": "The top-level src file",
        "doc": "The starlark_extract_doc output file",
        "deps": "DepSet[StarlarkLibraryFileInfo]: deps of this file",
        # "transitive_deps": "List[DepSet[StarlarkLibraryFileInfo]]: transitive deps of this file",
        "transitive_srcs": "Transitive closure of rules files required for " +
                           "interpretation of the src",
        "transitive_docs": "Transitive closure of docs that have viable dependencies",
        "broken": "If at last one of the transitive srcs has an unknown dependency.",
    },
)

def _starlark_library_impl(ctx):
    src = ctx.file.src
    doc = ctx.file.doc
    broken = len(ctx.attr.unknown_deps) > 0

    deps = [d[StarlarkLibraryFileInfo] for d in ctx.attr.deps]
    transitive_deps = [d.deps for d in deps]

    transitive_srcs = depset([src], order = "postorder", transitive = [d.transitive_srcs for d in deps])
    transitive_docs = depset([doc] if not broken else [], order = "postorder", transitive = [d.transitive_docs for d in deps])

    return [
        DefaultInfo(
            files = transitive_srcs,
        ),
        OutputGroupInfo(
            doc = [doc],
        ),
        StarlarkLibraryInfo(
            srcs = [src],
            transitive_srcs = transitive_srcs,
        ),
        StarlarkLibraryFileInfo(
            label = ctx.label,
            repo_name = ctx.attr.repo_name,
            src = src,
            deps = depset(deps, transitive = transitive_deps),
            transitive_srcs = transitive_srcs,
            transitive_docs = transitive_docs,
            broken = broken,
        ),
    ]

_starlark_library = rule(
    implementation = _starlark_library_impl,
    attrs = {
        "repo_name": attr.string(
            mandatory = True,
        ),
        "src": attr.label(
            doc = "label for the .bzl file",
            allow_single_file = True,
        ),
        "doc": attr.label(
            doc = "the output of the starlar_doc_extract rule",
            allow_single_file = True,
        ),
        "unknown_deps": attr.string_list(
            doc = "list of starlark_library rule dependencies that will not be able to resolve",
        ),
        "deps": attr.label_list(
            doc = "list of starlark_library rule dependencies.  These can be ",
            providers = [StarlarkLibraryFileInfo],
        ),
    },
    doc = "",
    provides = [DefaultInfo, StarlarkLibraryInfo, StarlarkLibraryFileInfo],
)

def starlark_library(name, src, deps = [], **kwargs):
    visibility = kwargs.pop("visibility", [])
    bazel_tools_deps = kwargs.pop("bazel_tools_deps", [])

    # doc_name = name + "_doc"
    # native.starlark_doc_extract(
    #     name = doc_name,
    #     src = src,
    #     deps = deps + bazel_tools_deps,
    #     visibility = ["//visibility:public"],
    # )

    _starlark_library(
        name = name,
        src = src,
        # doc = doc_name,
        doc = src,
        deps = deps,
        visibility = ["//visibility:public"],
        **kwargs
    )
