load("@build_stack_rules_proto//rules:providers.bzl", "ProtoRepositoryInfo")

def _proto_repository_info(ctx):
    return [
        ProtoRepositoryInfo(
            commit = ctx.attr.commit,
            tag = ctx.attr.tag,
            vcs = ctx.attr.vcs,
            urls = ctx.attr.urls,
            sha256 = ctx.attr.sha256,
            strip_prefix = ctx.attr.strip_prefix,
            source_host = ctx.attr.source_host,
            source_owner = ctx.attr.source_owner,
            source_repo = ctx.attr.source_repo,
            source_prefix = ctx.attr.source_prefix,
            source_commit = ctx.attr.source_commit,
        ),
    ]

proto_repository_info = rule(
    implementation = _proto_repository_info,
    attrs = {
        "commit": attr.string(doc = "the proto_repository.commit attr value"),
        "tag": attr.string(doc = "the proto_repository.tag attr value"),
        "vcs": attr.string(doc = "the proto_repository.vcs attr value"),
        "urls": attr.string_list(doc = "the proto_repository.urls attr value"),
        "sha256": attr.string(doc = "the proto_repository.sha256 attr value"),
        "strip_prefix": attr.string(doc = "the proto_repository.strip_prefix attr value"),
        "source_host": attr.string(doc = "the proto_repository.source_host attr value"),
        "source_owner": attr.string(doc = "the proto_repository.source_owner attr value"),
        "source_repo": attr.string(doc = "the proto_repository.source_repo attr value"),
        "source_prefix": attr.string(doc = "the proto_repository.source_prefix attr value"),
        "source_commit": attr.string(doc = "the proto_repository.source_commit attr value"),
    },
)
