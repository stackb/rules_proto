"mypy_stubs rule"
MyPyStubsInfo = provider(
    "MyPyStubsInfo provider",
    fields = {
        "srcs": ".pyi stub files",
    },
)

def _mypy_stubs_impl(ctx):
    pyi_srcs = []
    for target in ctx.attr.srcs:
        pyi_srcs.extend(target.files.to_list())
    transitive_srcs = depset(direct = pyi_srcs)

    return [
        MyPyStubsInfo(
            srcs = ctx.attr.srcs,
        ),
        PyInfo(
            has_py2_only_sources = False,
            has_py3_only_sources = True,
            uses_shared_libraries = False,
            transitive_sources = transitive_srcs,
        ),
    ]

mypy_stubs = rule(
    implementation = _mypy_stubs_impl,
    attrs = {
        "srcs": attr.label_list(
            allow_empty = False,
            mandatory = True,
            allow_files = [".pyi"],
        ),
    },
)
