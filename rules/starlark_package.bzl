"""starlark_package.bzl is similar to bzl_library but for a BUILD.package file."""

StarlarkPackageInfo = provider(
    "Information about a single .bzl file.",
    fields = {
        "label": "Label: The label of the target rule",
        "src": "File: The .bzl file",
    },
)

def _starlark_package_impl(ctx):
    return [
        StarlarkPackageInfo(
            label = ctx.label,
            src = ctx.file.src,
        ),
    ]

starlark_package = rule(
    implementation = _starlark_package_impl,
    attrs = {
        "src": attr.label(
            doc = "the .package source file",
            allow_single_file = True,
        ),
    },
    provides = [StarlarkPackageInfo],
)
