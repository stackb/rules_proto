"""protogenrule wraps native.genrule with supplemental .update and _test targets

protogenrule is a drop-in replacement for native.genrule that provides two additional targets
used to copying source back into the monorepo and asserting that the
generated file(s) remain consistent with the source control version of the file.

For any protogenrule target `//:a` that outputs `a.txt`, `//:a.update`
copies `a.txt` back into the source tree as `a.txt'`; `//:a_test` asserts that
`a.txt` and `a.txt'` are identical.
"""

load(
    "@build_stack_rules_proto//rules:proto_compile_gencopy.bzl",
    "proto_compile_gencopy_run",
    "proto_compile_gencopy_test",
)
load(
    "@build_stack_rules_proto//rules:providers.bzl",
    "ProtoCompileInfo",
)

def _proto_compiled_sources_impl(ctx):
    dep = ctx.attr.dep[DefaultInfo]
    return ProtoCompileInfo(
        label = ctx.attr.dep.label,
        outputs = dep.files.to_list(),
    )

_proto_compiled_sources = rule(
    doc = """Provider Adapter from DefaultInfo to ProtoCompileInfo.
        """,
    implementation = _proto_compiled_sources_impl,
    attrs = {"dep": attr.label(providers = [DefaultInfo])},
)

def protogenrule(
        name,
        run_target_suffix = ".update",
        sources_target_suffix = "d_sources",
        test_target_suffix = "_test",
        **kwargs):
    """protogenrule is used identically to native.gencopy

    Args:
        name: the name of the rule
        run_target_suffix: the suffix for the update/copy target
        sources_target_suffix: the suffix for the _proto_compiled_sources target
        test_target_suffix: the suffix for the test target
        **kwargs: remainder of non-positional args
    """
    name_run = name + run_target_suffix
    name_sources = name + sources_target_suffix
    name_test = name + test_target_suffix

    outs = kwargs.pop("outs", [])
    visibility = kwargs.pop("visibility", [])

    native.genrule(
        name = name,
        outs = outs,
        visibility = visibility,
        **kwargs
    )

    _proto_compiled_sources(
        name = name_sources,
        dep = name,
        visibility = visibility,
    )

    proto_compile_gencopy_test(
        name = name_test,
        srcs = outs,
        deps = [name_sources],
        mode = "check",
        update_target_label_name = name_run,
    )

    proto_compile_gencopy_run(
        name = name_run,
        deps = [name_sources],
        mode = "update",
        update_target_label_name = name_run,
    )
