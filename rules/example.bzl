"example.bzl provides the gazelle_testdata_example rule."

load("@io_bazel_rules_go//go/tools/bazel_testing:def.bzl", "go_bazel_test")

def _examplegen_impl(ctx):
    config_json = ctx.outputs.json
    output_test = ctx.outputs.test
    output_markdown = ctx.outputs.markdown

    config = struct(
        name = ctx.label.name,
        label = str(ctx.label),
        testOut = output_test.path,
        markdownOut = output_markdown.path,
        workspaceIn = ctx.file.workspace_template.path,
        stripPrefix = ctx.attr.strip_prefix,
        files = [f.path for f in ctx.files.srcs],
    )

    ctx.actions.write(
        output = config_json,
        content = config.to_json(),
    )

    ctx.actions.run(
        mnemonic = "ExampleGenerate",
        progress_message = "Generating %s test" % ctx.attr.name,
        executable = ctx.file._examplegen,
        arguments = ["--config_json=%s" % config_json.path],
        inputs = [config_json, ctx.file.workspace_template] + ctx.files.srcs,
        outputs = [output_test, output_markdown],
    )

    return [DefaultInfo(
        files = depset([config_json, output_test, output_markdown]),
    )]

_examplegen = rule(
    implementation = _examplegen_impl,
    attrs = {
        "srcs": attr.label_list(
            doc = "Sources for the test txtar file",
            allow_files = True,
        ),
        "strip_prefix": attr.string(
            doc = "path prefix to remove from test files in the txtar",
        ),
        "workspace_template": attr.label(
            doc = "Template for the test WORKSPACE",
            allow_single_file = True,
            mandatory = True,
        ),
        "_examplegen": attr.label(
            doc = "The examplegen generator tool",
            default = "//cmd/examplegen",
            allow_single_file = True,
            executable = True,
            cfg = "host",
        ),
    },
    outputs = {
        "json": "%{name}.json",
        "test": "%{name}_test.go",
        "markdown": "%{name}.md",
    },
)

def gazelle_testdata_example(**kwargs):
    """
    gazelle_testdata_example rule runs an go_bazel_test for an example dir

    Args:
        **kwargs: the kwargs dict passed to 'go_bazel_test'
    """
    name = kwargs.pop("name")
    srcs = kwargs.pop("srcs", [])
    strip_prefix = kwargs.pop("strip_prefix", "")

    rule_files = kwargs.pop("rule_files", ["//:all_files"])

    _examplegen(
        name = name,
        srcs = srcs,
        strip_prefix = strip_prefix,
        workspace_template = kwargs.pop("workspace_template", ""),
    )

    go_bazel_test(
        name = name + "_test",
        srcs = [name + "_test.go"],
        rule_files = rule_files,
        **kwargs
    )
