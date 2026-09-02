"""starlark_package_library.bzl is similar to bzl_library but also provides load statement foreach file."""

load("//rules:starlark_package.bzl", "StarlarkPackageInfo")

StarlarkPackageLibraryInfo = provider(
    "Information on a set of starlark packages.  This is a flat list, non-transitive.",
    fields = {
        "label": "The label of the target rule",
        "packages": "List[StarlarkPackageInfo]: package deps of this rule",
        "srcs": "List[File]: source files for the packages, for convenience",
        "bazelignore": "List[str] value of ctx.attr.bazelignore",
        "bazelversion": "str: the value of ctx.attr.bazelversion",
    },
)

def _starlark_package_library_impl(ctx):
    packages = [m[StarlarkPackageInfo] for m in ctx.attr.packages]
    return [
        StarlarkPackageLibraryInfo(
            label = ctx.label,
            bazelignore = ctx.attr.bazelignore,
            bazelversion = ctx.attr.bazelversion,
            packages = packages,
            srcs = [m.src for m in packages],
        ),
    ]

starlark_package_library = rule(
    implementation = _starlark_package_library_impl,
    attrs = {
        "bazelignore": attr.string_list(
            doc = "contents of the .bazelignore file, if present",
        ),
        "bazelversion": attr.string(
            doc = "contents of the .bazelversion file, if present",
        ),
        "packages": attr.label_list(
            doc = "list of starlark_package rule dependencies.",
            providers = [StarlarkPackageInfo],
        ),
    },
    provides = [StarlarkPackageLibraryInfo],
)
