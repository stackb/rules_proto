"""starlark_module.bzl is similar to bzl_library but also provides load statement foreach file."""

load("@bazel_skylib//:bzl_library.bzl", "StarlarkLibraryInfo")

StarlarkModuleInfo = provider(
    "Information about a single .bzl file.",
    fields = {
        "label": "Label: The label of the target rule",
        "loads": "List[str]: load statements for this file",
        "src": "File: The .bzl file",
    },
)

def _starlark_module_impl(ctx):
    return [
        StarlarkLibraryInfo(
            srcs = [ctx.file.src],
            transitive_srcs = depset([ctx.file.src]),
        ),
        StarlarkModuleInfo(
            label = ctx.label,
            loads = ctx.attr.loads,
            src = ctx.file.src,
        ),
    ]

starlark_module = rule(
    implementation = _starlark_module_impl,
    attrs = {
        "loads": attr.string_list(
            doc = "the load statements in this file",
        ),
        "src": attr.label(
            doc = "the .bzl source file",
            allow_single_file = True,
        ),
    },
    provides = [StarlarkLibraryInfo, StarlarkModuleInfo],
)
