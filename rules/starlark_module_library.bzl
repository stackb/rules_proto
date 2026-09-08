"""starlark_module_library.bzl is similar to bzl_library but also provides load statement foreach file."""

load("//rules:starlark_module.bzl", "StarlarkModuleInfo")

StarlarkModuleLibraryInfo = provider(
    "Information on a set of starlark modules.  This is a flat list, non-transitive.",
    fields = {
        "label": "The label of the target rule",
        "modules": "List[StarlarkModuleInfo]: modules deps of this rule",
        "srcs": "List[File]: source files for the modules, for convenience",
        "bazelignore": "List[str] value of ctx.attr.bazelignore",
        "bazelversion": "str: the value of ctx.attr.bazelversion",
    },
)

def _starlark_module_library_impl(ctx):
    modules = [m[StarlarkModuleInfo] for m in ctx.attr.modules]
    return [
        StarlarkModuleLibraryInfo(
            label = ctx.label,
            bazelignore = ctx.attr.bazelignore,
            bazelversion = ctx.attr.bazelversion,
            modules = modules,
            srcs = [m.src for m in modules],
        ),
    ]

starlark_module_library = rule(
    implementation = _starlark_module_library_impl,
    attrs = {
        "bazelignore": attr.string_list(
            doc = "contents of the .bazelignore file, if present",
        ),
        "bazelversion": attr.string(
            doc = "contents of the .bazelversion file, if present",
        ),
        "modules": attr.label_list(
            doc = "list of starlark_module rule dependencies.",
            providers = [StarlarkModuleInfo],
        ),
    },
    provides = [StarlarkModuleLibraryInfo],
)
