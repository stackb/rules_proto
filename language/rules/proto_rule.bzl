load(
    "@build_stack_rules_proto//rules:proto_plugin.bzl",
    "ProtoPluginInfo",
    "proto_plugin_info_to_struct",
)
load(
    "//language/rules:proto_language.bzl",
    "ProtoLanguageInfo",
    "proto_language_info_to_struct",
)
load(
    "@build_stack_rules_proto//:provider_test.bzl",
    "redact_host_configuration",
)

ProtoRuleInfo = provider("Provider for a proto rule", fields = {
    "name": "The prefix name of the rule (e.g 'py')",
    "rule": "The rule struct",
    "bzl_file": "The generated rule.bzl file",
    "build_file": "The generated example BUILD file",
    "workspace_file": "The generated example WORKSPACE file",
    "markdown_file": "The generated markdown/README file",
    "deps_file": "The generated deps file",
    "bazel_test_file": "The generated bazel_test file",
    "json_file": "The generated json file",
})

def proto_rule_info_to_struct(info):
    return struct(
        name = info.name,
        rule = rule_to_struct(info.rule),
        bzl_file = redact_host_configuration(info.bzl_file.short_path),
        build_file = redact_host_configuration(info.build_file.short_path),
        workspace_file = redact_host_configuration(info.workspace_file.short_path),
        deps_file = redact_host_configuration(info.deps_file.short_path),
        bazel_test_file = redact_host_configuration(info.bazel_test_file.short_path),
        json_file = redact_host_configuration(info.json_file.short_path),
    )

def rule_to_struct(rule):
    return struct(
        name = rule.name,
        package = rule.package,
        skipDirectoriesMerge = rule.skipDirectoriesMerge,
        implementationFilename = redact_host_configuration(rule.implementationFilename),
        implementationTmpl = redact_host_configuration(rule.implementationTmpl),
        workspaceExampleFilename = redact_host_configuration(rule.workspaceExampleFilename),
        workspaceExampleTmpl = redact_host_configuration(rule.workspaceExampleTmpl),
        buildExampleFilename = redact_host_configuration(rule.buildExampleFilename),
        buildExampleTmpl = redact_host_configuration(rule.buildExampleTmpl),
        markdownFilename = redact_host_configuration(rule.markdownFilename),
        markdownTmpl = redact_host_configuration(rule.markdownTmpl),
        depsFilename = redact_host_configuration(rule.depsFilename),
        depsTmpl = redact_host_configuration(rule.depsTmpl),
        bazelTestFilename = redact_host_configuration(rule.bazelTestFilename),
        bazelTestTmpl = redact_host_configuration(rule.bazelTestTmpl),
        plugins = rule.plugins,
    )

def _proto_rule_impl(ctx):
    rule_json = ctx.outputs.json
    output_bzl = ctx.outputs.bzl
    output_workspace = ctx.outputs.workspace
    output_build = ctx.outputs.build
    output_markdown = ctx.outputs.markdown
    output_deps = ctx.outputs.deps
    output_bazel_test = ctx.outputs.bazel_test

    rule = struct(
        name = ctx.attr.name,
        package = ctx.attr.package or ctx.label.package,
        skipDirectoriesMerge = ctx.attr.skip_directories_merge,
        implementationFilename = output_bzl.path,
        implementationTmpl = ctx.file.implementation_tmpl.path,
        workspaceExampleFilename = output_workspace.path,
        workspaceExampleTmpl = ctx.file.workspace_example_tmpl.path,
        buildExampleFilename = output_build.path,
        buildExampleTmpl = ctx.file.build_example_tmpl.path,
        markdownFilename = output_markdown.path,
        markdownTmpl = ctx.file.markdown_tmpl.path,
        depsFilename = output_deps.path,
        depsTmpl = ctx.file.deps_tmpl.path,
        bazelTestFilename = output_bazel_test.path,
        bazelTestTmpl = ctx.file.bazel_test_tmpl.path,
        plugins = [proto_plugin_info_to_struct(p[ProtoPluginInfo]) for p in ctx.attr.plugins],
        language = proto_language_info_to_struct(ctx.attr.language[ProtoLanguageInfo]),
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
        ctx.file.deps_tmpl,
        ctx.file.bazel_test_tmpl,
    ]

    outputs = [
        output_bzl,
        output_build,
        output_workspace,
        output_markdown,
        output_deps,
        output_bazel_test,
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
            deps_file = output_deps,
            bazel_test_file = output_bazel_test,
            json_file = rule_json,
        ),
        DefaultInfo(
            files = depset(outputs + [rule_json]),
        ),
    ]

proto_rule = rule(
    implementation = _proto_rule_impl,
    attrs = {
        "package": attr.string(
            doc = "The target package for the rule. If empty, default to ctx.label.package",
        ),
        "implementation_tmpl": attr.label(
            doc = "The rule implementation template",
            default = str(Label("//tools/protogen:aspect.bzl.tmpl")),
            allow_single_file = True,
        ),
        "workspace_example_tmpl": attr.label(
            doc = "The rule workspace example template",
            default = str(Label("//tools/protogen:WORKSPACE.tmpl")),
            allow_single_file = True,
        ),
        "build_example_tmpl": attr.label(
            doc = "The rule build example template",
            default = str(Label("//tools/protogen:BUILD.tmpl")),
            allow_single_file = True,
        ),
        "markdown_tmpl": attr.label(
            doc = "The rule build markdown example template",
            default = str(Label("//tools/protogen:proto_rule.md.tmpl")),
            allow_single_file = True,
        ),
        "deps_tmpl": attr.label(
            doc = "The workspace deps example template",
            default = str(Label("//tools/protogen:deps.bzl.tmpl")),
            allow_single_file = True,
        ),
        "bazel_test_tmpl": attr.label(
            doc = "The workspace bazel_test example template",
            default = str(Label("//tools/protogen:bazel_test.go.tmpl")),
            allow_single_file = True,
        ),
        "plugins": attr.label_list(
            doc = "List of default plugins to include in the generated rule",
            providers = [ProtoPluginInfo],
        ),
        "language": attr.label(
            doc = "The language this rule belongs to",
            providers = [ProtoLanguageInfo],
            mandatory = True,
        ),
        "skip_directories_merge": attr.bool(
            doc = "If the generated rule shoul skip merging directories",
        ),
        "data": attr.label_list(
            allow_files = True,
        ),
        "_rulegen": attr.label(
            doc = "The rulegen generator tool",
            default = "//tools/protogen/cmd/rulegen",
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
        "deps": "%{name}_deps.bzl",
        "bazel_test": "%{name}_bazel_test.go",
    },
)
