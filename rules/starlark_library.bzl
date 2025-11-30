"""starlark_library.bzl is a thin wrapper over bzl_library."""

load("@bazel_skylib//:bzl_library.bzl", "StarlarkLibraryInfo")

StarlarkLibraryFileInfo = provider(
    "Information on contained Starlark rules.",
    fields = {
        "label": "The label of the target rule",
        "repo_name": "The original (non-canonical) repo (workspace!) name",
        "srcs": "List of .bzl files",
        "loads": "load statements per .bzl file",
        "deps": "DepSet[StarlarkLibraryFileInfo]: deps of this file",
        "transitive_deps": "DepSet[StarlarkLibraryFileInfo]: transitive deps of this file",
        "bazelignore": "List[str] value of ctx.attr.bazelignore",
        "bazelversion": "str: the value of ctx.attr.bazelversion",
    },
)

def _starlark_library_impl(ctx):
    deps = [d[StarlarkLibraryFileInfo] for d in ctx.attr.deps]
    srcs = ctx.files.srcs
    transitive_srcs = depset(srcs, order = "postorder", transitive = [d.transitive_srcs for d in deps])

    # loads = {}
    # for k, v in ctx.attr.loads.items():
    #     loads[Label(k)] = v

    return [
        DefaultInfo(
            files = transitive_srcs,
        ),
        StarlarkLibraryInfo(
            srcs = srcs,
            transitive_srcs = transitive_srcs,
        ),
        StarlarkLibraryFileInfo(
            label = ctx.label,
            repo_name = ctx.attr.repo_name,
            srcs = srcs,
            loads = ctx.attr.loads,
            bazelignore = ctx.attr.bazelignore,
            bazelversion = ctx.attr.bazelversion,
            transitive_deps = depset(deps, transitive = [d.transitive_deps for d in deps]),
        ),
    ]

starlark_library = rule(
    implementation = _starlark_library_impl,
    attrs = {
        "repo_name": attr.string(
            mandatory = True,
        ),
        "srcs": attr.label_list(
            doc = "label for the .bzl file",
            allow_files = True,
        ),
        "loads": attr.string_list_dict(
            doc = "load per file",
        ),
        "deps": attr.label_list(
            doc = "list of starlark_library rule dependencies.  These can be ",
            providers = [StarlarkLibraryFileInfo],
        ),
        "bazelignore": attr.string_list(
            doc = "contents of the .bazelignore file, if present",
        ),
        "bazelversion": attr.string(
            doc = "contents of the .bazelversion file, if present",
        ),
    },
    doc = "",
    provides = [DefaultInfo, StarlarkLibraryInfo, StarlarkLibraryFileInfo],
)
