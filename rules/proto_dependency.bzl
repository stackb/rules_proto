"""proto_dependency.bzl 
"""
load("//rules:providers.bzl", "ProtoDependencyInfo")

def _proto_dependency_impl(ctx):
    return [
        ProtoDependencyInfo(
            buildFile = ctx.attr.build_file,
            deps = [dep[ProtoDependencyInfo] for dep in ctx.attr.deps],
            label = str(ctx.label),
            name = ctx.attr.name,
            repositoryRule = ctx.attr.repository_rule,
            sha256 = ctx.attr.sha256,
            stripPrefix = ctx.attr.strip_prefix,
            urls = ctx.attr.urls,
            workspaceSnippet = ctx.attr.workspace_snippet,
        ),
    ]

proto_dependency = rule(
    implementation = _proto_dependency_impl,
    attrs = {
        "build_file": attr.string(
            doc = "The build_file attribute for http_archive",
        ),
        "workspace_snippet": attr.string(
            doc = "The starlark code snippet for the WORKSPACE needed when using this dependency",
        ),
        "deps": attr.label_list(
            doc = "Additional transitive dependencies",
            providers = [ProtoDependencyInfo],
        ),
        "repository_rule": attr.string(
            doc = "The repository rule that instantiates this dependency",
            values = ["http_archive", "http_file", "bind", "go_repository", "phony"],
        ),
        "sha256": attr.string(
            doc = "The sha256 attribute for http_archive",
        ),
        "strip_prefix": attr.string(
            doc = "The strip_prefix attribute for http_archive",
        ),
        "urls": attr.string_list(
            doc = "The strip_prefix attribute for http_archive",
        ),
    },
)