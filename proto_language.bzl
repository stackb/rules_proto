load(
    "//:proto_dependency.bzl",
    "ProtoDependencyInfo",
    "proto_dependency_info_to_struct",
)
load(
    "@build_stack_rules_proto//:provider_test.bzl",
    "redact_host_configuration",
)

ProtoLanguageInfo = provider("Provider for a proto language", fields = {
    "name": "The name of the language (e.g 'python')",
    "label": "The label of the language",
    "rules": "The list of rules that this language is associated with",
    "markdown_file": "The markdown documentation file for this language",
    "deps": "List of ProtoDependencyInfo",
})

def proto_language_info_to_struct(info):
    return struct(
        name = info.name,
        label = str(info.label),
        rules = info.rules,
        markdown_file = redact_host_configuration(info.markdown_file.short_path),
        deps = [proto_dependency_info_to_struct(dep) for dep in info.deps],
    )

def _proto_language_impl(ctx):
    return [
        ProtoLanguageInfo(
            name = ctx.attr.name,
            label = ctx.label,
            rules = ctx.attr.rules,
            markdown_file = ctx.file.markdown_tmpl,
            deps = depset(
                direct = [dep[ProtoDependencyInfo] for dep in ctx.attr.deps],
                transitive = [dep[ProtoDependencyInfo].deps for dep in ctx.attr.deps],
            ),
        ),
    ]

proto_language = rule(
    implementation = _proto_language_impl,
    attrs = {
        "rules": attr.string_list(
            doc = "The list of rules that belong to this language",
        ),
        "markdown_tmpl": attr.label(
            doc = "The rule build markdown example template",
            default = str(Label("//tools/protorule:language.md.tmpl")),
            allow_single_file = True,
        ),
        "deps": attr.label_list(
            doc = "List of deps that apply to all rules belonging to this language",
            providers = [ProtoDependencyInfo],
        ),
    },
    outputs = {
        "markdown": "%{name}.md",
    },
)
