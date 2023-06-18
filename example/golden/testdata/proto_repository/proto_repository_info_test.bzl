load("@build_stack_rules_proto//rules:providers.bzl", "ProtoRepositoryInfo")

def _proto_repository_info_test_impl(ctx):
    info = ctx.attr.info[ProtoRepositoryInfo]

    must_attr(info, ctx.attr, "source_host")
    must_attr(info, ctx.attr, "source_owner")
    must_attr(info, ctx.attr, "source_repo")
    must_attr(info, ctx.attr, "source_commit")
    must_attr(info, ctx.attr, "source_prefix")

    ctx.actions.write(ctx.outputs.json, info.to_json())

    # we're checking attr values in the provider, so the script really does not
    # need to do anything
    ctx.actions.write(ctx.outputs.executable, "echo PASS")

    return [DefaultInfo(
        files = depset([ctx.outputs.json, ctx.outputs.executable]),
    )]

proto_repository_info_test = rule(
    implementation = _proto_repository_info_test_impl,
    attrs = {
        "info": attr.label(
            providers = [ProtoRepositoryInfo],
            mandatory = True,
        ),
        "want_source_host": attr.string(),
        "want_source_owner": attr.string(),
        "want_source_repo": attr.string(),
        "want_source_commit": attr.string(),
        "want_source_prefix": attr.string(),
    },
    outputs = {
        "json": "%{name}.json",
    },
    test = True,
)

def must_attr(info, attr, name):
    got = getattr(info, name)
    want = getattr(attr, "want_" + name)
    if got != want:
        fail(".%s: want %s, got %s" % (name, want, got))
