"""golden_filegroup wraps native.filegroup with .update and .test targets

golden_filegroup is a drop-in replacement for native.filegroup that provides
two additional targets used to copying source back into the source tree and
asserting that the file(s) remain consistent with the source control version
of the file (the golden files).

For any golden_filegroup target `//:a` that srcs `a.txt`, `//:a.update` copies
`a.txt` back into the source tree as `a.txt.golden'`; `//:a.test` asserts that
`a.txt` and `a.txt.golden'` are identical.
"""

load("@build_stack_rules_proto//rules:providers.bzl", "ProtoCompileInfo")
load(
    "@build_stack_rules_proto//rules:proto_compile_gencopy.bzl",
    "proto_compile_gencopy_run",
    "proto_compile_gencopy_test",
)

def _files_impl(ctx):
    dep = ctx.attr.dep[DefaultInfo]
    return ProtoCompileInfo(
        label = ctx.attr.dep.label,
        outputs = dep.files.to_list(),
    )

_files = rule(
    doc = """Provider Adapter from DefaultInfo to ProtoCompileInfo.""",
    implementation = _files_impl,
    attrs = {"dep": attr.label(providers = [DefaultInfo])},
)

def golden_filegroup(
        name,
        sources_target_suffix = ".files",
        test_target_suffix = ".test",
        run_target_suffix = ".update",
        extension = ".golden",
        **kwargs):
    """golden_filegroup is used identically to native.gencopy

    Args:
        name: the name of the rule
        run_target_suffix: the suffix for the update/copy target
        sources_target_suffix: the suffix for the _proto_compiled_sources target
        test_target_suffix: the suffix for the test target
        extension: the golden file extension to append
        **kwargs: remainder of non-positional args
    """
    name_sources = name + sources_target_suffix
    name_test = name + test_target_suffix
    name_run = name + run_target_suffix

    tags = kwargs.pop("tags", [])
    srcs = kwargs.pop("srcs", [])
    goldens = [src + extension for src in srcs]
    native.filegroup(name = name, srcs = srcs, tags = tags, **kwargs)

    _files(name = name_sources, dep = name, tags = tags)

    proto_compile_gencopy_test(
        name = name_test,
        srcs = goldens,
        deps = [name_sources],
        mode = "check",
        extension = extension,
        update_target_label_name = name_run,
        tags = tags,
    )

    proto_compile_gencopy_run(
        name = name_run,
        deps = [name_sources],
        mode = "update",
        extension = extension,
        update_target_label_name = name_run,
        tags = tags,
    )
