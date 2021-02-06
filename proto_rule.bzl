load("//:plugin.bzl", "ProtoPluginInfo")

ProtoRuleInfo = provider("Provider for a proto rule", fields = {
    "name": "The prefix name of the rule (e.g 'py')",
    "rule": "The rule struct",
    "bzl_file": "The generated rule.bzl file",
    "build_file": "The generated example BUILD file",
    "workspace_file": "The generated example WORKSPACE file",
    "markdown_file": "The generated markdown/README file",
})

def _proto_rule_impl(ctx):
    rule_json = ctx.outputs.json
    output_bzl = ctx.outputs.bzl
    output_workspace = ctx.outputs.workspace
    output_build = ctx.outputs.build
    output_markdown = ctx.outputs.markdown

    rule = struct(
        name = ctx.attr.name,
        kind = ctx.attr.kind,
        package = ctx.label.package,
        skipDirectoriesMerge = ctx.attr.skip_directories_merge,
        plugins = [str(p.label) for p in ctx.attr.plugins],
        implementationFilename = output_bzl.path,
        implementationTmpl = ctx.file.implementation_tmpl.path,
        workspaceExampleFilename = output_workspace.path,
        workspaceExampleTmpl = ctx.file.workspace_example_tmpl.path,
        buildExampleFilename = output_build.path,
        buildExampleTmpl = ctx.file.build_example_tmpl.path,
        markdownFilename = output_markdown.path,
        markdownTmpl = ctx.file.markdown_tmpl.path,
    )

    ctx.actions.write(
        output = rule_json,
        content = rule.to_json(),
    )

    inputs = [
        rule_json,
        ctx.file.implementation_tmpl,
        ctx.file.build_example_tmpl,
        ctx.file.workspace_example_tmpl,
        ctx.file.markdown_tmpl,
    ]

    outputs = [
        output_bzl,
        output_build,
        output_workspace,
        output_markdown,
    ]

    args = [
        "--rule_json=%s" % rule_json.path,
    ]

    ctx.actions.run(
        mnemonic = "ProtoRuleGenerate",
        progress_message = "Generating %s rule" % ctx.attr.name,
        executable = ctx.file._rulegen,
        arguments = args,
        inputs = inputs,
        outputs = outputs,
    )

    return [
        ProtoRuleInfo(
            name = ctx.attr.name,
            rule = rule,
            bzl_file = output_bzl,
            build_file = output_build,
            workspace_file = output_workspace,
            markdown_file = output_markdown,
        ),
        DefaultInfo(
            files = depset(outputs + [rule_json]),
        ),
    ]

proto_rule = rule(
    implementation = _proto_rule_impl,
    attrs = {
        "kind": attr.string(
            doc = "The kind of rule",
            values = ["proto", "grpc"],
        ),
        "implementation_tmpl": attr.label(
            doc = "The rule implementation template",
            default = str(Label("//tools/protorule:aspect.bzl.tmpl")),
            allow_single_file = True,
        ),
        "workspace_example_tmpl": attr.label(
            doc = "The rule workspace example template",
            default = str(Label("//tools/protorule:WORKSPACE.tmpl")),
            allow_single_file = True,
        ),
        "build_example_tmpl": attr.label(
            doc = "The rule build example template",
            default = str(Label("//tools/protorule:BUILD.tmpl")),
            allow_single_file = True,
        ),
        "markdown_tmpl": attr.label(
            doc = "The rule build markdown example template",
            default = str(Label("//tools/protorule:README.md.tmpl")),
            allow_single_file = True,
        ),
        "plugins": attr.label_list(
            doc = "List of default plugins to include in the generated rule",
            providers = [ProtoPluginInfo],
        ),
        "skip_directories_merge": attr.bool(
            doc = "If the generated rule shoul skip merging directories",
        ),
        "data": attr.label_list(allow_files = True),
        "_rulegen": attr.label(
            doc = "The rulegen generator tool",
            default = "//tools/protorule/cmd/rulegen",
            allow_single_file = True,
            executable = True,
            cfg = "host",
        ),
    },
    outputs = {
        "bzl": "%{name}.bzl",
        "markdown": "%{name}.md",
        "json": "%{name}.json",
        "workspace": "%{name}.WORKSPACE",
        "build": "%{name}.BUILD",
    },
)
