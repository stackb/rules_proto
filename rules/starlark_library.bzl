"""starlark_library.bzl is similar to bzl_library but also provides load statement foreach file."""

load("@bazel_skylib//:bzl_library.bzl", "StarlarkLibraryInfo")

StarlarkLibraryFileInfo = provider(
    "Information on contained Starlark rules.",
    fields = {
        "bazelignore": "List[str] value of ctx.attr.bazelignore",
        "bazelversion": "str: the value of ctx.attr.bazelversion",
        "deps": "DepSet[StarlarkLibraryFileInfo]: deps of this file",
        "label": "The label of the target rule",
        "loads": "load statements per .bzl file",
        "srcs": "List of .bzl files",
        "transitive_deps": "DepSet[StarlarkLibraryFileInfo]: transitive deps of this file",
    },
)

def _starlark_library_impl(ctx):
    srcs = ctx.files.srcs
    deps = [d[StarlarkLibraryFileInfo] for d in ctx.attr.deps]
    transitive_srcs = depset(srcs, order = "postorder", transitive = [d.transitive_srcs for d in deps])
    transitive_deps = depset(deps, order = "postorder", transitive = [d.transitive_deps for d in deps])

    return [
        DefaultInfo(
            files = transitive_srcs,
        ),
        StarlarkLibraryInfo(
            srcs = srcs,
            transitive_srcs = transitive_srcs,
        ),
        StarlarkLibraryFileInfo(
            bazelignore = ctx.attr.bazelignore,
            bazelversion = ctx.attr.bazelversion,
            deps = deps,
            label = ctx.label,
            loads = ctx.attr.loads,
            srcs = srcs,
            transitive_deps = transitive_deps,
        ),
    ]

starlark_library = rule(
    implementation = _starlark_library_impl,
    attrs = {
        "bazelignore": attr.string_list(
            doc = "contents of the .bazelignore file, if present",
        ),
        "bazelversion": attr.string(
            doc = "contents of the .bazelversion file, if present",
        ),
        "loads": attr.string_list_dict(
            doc = "load per file",
        ),
        "srcs": attr.label_list(
            doc = "label for the .bzl file",
            allow_files = True,
        ),
        "deps": attr.label_list(
            doc = "list of starlark_library rule dependencies.  These can be ",
            providers = [StarlarkLibraryFileInfo],
        ),
    },
    doc = "",
    provides = [DefaultInfo, StarlarkLibraryInfo, StarlarkLibraryFileInfo],
)
