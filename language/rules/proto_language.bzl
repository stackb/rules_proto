load(
    "@build_stack_rules_proto//:provider_test.bzl",
    "redact_host_configuration",
)

ProtoLanguageInfo = provider("Provider for a proto language", fields = {
    "name": "The name of the language (e.g 'python')",
    "prefix": "The language prefix (e.g 'py')",
    "lang": "The lang struct (config file for langgen)",
    "markdown_file": "The markdown documentation file for this language",
    "rules_file": "The rules bzl file for this language",
})

def proto_language_info_to_struct(info):
    return struct(
        lang = lang_to_struct(info.lang),
        markdown_file = redact_host_configuration(info.markdown_file.short_path),
        name = info.name,
        prefix = info.prefix,
        rules_file = redact_host_configuration(info.rules_file.short_path),
    )

def lang_to_struct(lang):
    return struct(
        markdownFilename = redact_host_configuration(lang.markdownFilename),
        markdownTmpl = redact_host_configuration(lang.markdownTmpl),
        name = lang.name,
        package = lang.package,
        prefix = lang.prefix,
        rules = lang.rules,
        rulesFilename = redact_host_configuration(lang.rulesFilename),
        rulesTmpl = redact_host_configuration(lang.rulesTmpl),
    )

def _proto_language_impl(ctx):
    language_json = ctx.outputs.json
    output_markdown = ctx.outputs.markdown
    output_rules = ctx.outputs.rules

    lang = struct(
        markdownFilename = output_markdown.path,
        markdownTmpl = ctx.file.markdown_tmpl.path,
        name = ctx.attr.name,
        package = ctx.attr.package or ctx.label.package,
        prefix = ctx.attr.prefix,
        rules = ctx.attr.rules,
        rulesFilename = output_rules.path,
        rulesTmpl = ctx.file.rules_tmpl.path,
    )

    ctx.actions.write(
        output = language_json,
        content = lang.to_json(),
    )

    inputs = [
        language_json,
        ctx.file.markdown_tmpl,
        ctx.file.rules_tmpl,
    ]

    outputs = [
        output_markdown,
        output_rules,
    ]

    args = [
        "--language_json=%s" % language_json.path,
    ]

    ctx.actions.run(
        mnemonic = "ProtoLanguageGenerate",
        progress_message = "Generating %s language" % ctx.attr.name,
        executable = ctx.file._langgen,
        arguments = args,
        inputs = inputs,
        outputs = outputs,
    )

    return [
        ProtoLanguageInfo(
            name = ctx.attr.name,
            prefix = ctx.attr.prefix,
            lang = lang,
            markdown_file = output_markdown,
            rules_file = output_rules,
        ),
        DefaultInfo(
            files = depset(outputs + [language_json]),
        ),
    ]

proto_language = rule(
    implementation = _proto_language_impl,
    attrs = {
        "prefix": attr.string(
            doc = "The language prefix",
            mandatory = True,
        ),
        "package": attr.string(
            doc = "The target package for the language generated outputs. If empty, default to ctx.label.package",
        ),
        "rules": attr.string_list(
            doc = "The list of rules that belong to this language",
        ),
        "rules_tmpl": attr.label(
            doc = "The language rules example template",
            default = str(Label("//tools/protogen:rules.bzl.tmpl")),
            allow_single_file = True,
        ),
        "markdown_tmpl": attr.label(
            doc = "The rule build markdown example template",
            default = str(Label("//tools/protogen:proto_language.md.tmpl")),
            allow_single_file = True,
        ),
        "_langgen": attr.label(
            doc = "The langgen generator tool",
            default = "//tools/protogen/cmd/langgen",
            allow_single_file = True,
            executable = True,
            cfg = "host",
        ),
    },
    outputs = {
        "markdown": "%{name}.md",
        "rules": "%{name}.bzl",
        "json": "%{name}.json",
    },
)
