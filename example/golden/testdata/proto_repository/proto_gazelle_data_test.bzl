"""proto_gazelle_data_test is a regression test.

See https://github.com/stackb/rules_proto/issues/342.

"""

def _proto_gazelle_data_test_impl(ctx):
    info = ctx.attr.gazelle[DefaultInfo]

    ctx.actions.write(ctx.outputs.json, struct(
        data_files = [f.short_path for f in info.files.to_list()],
    ).to_json())

    # we're checking attr values in the provider, so the script really does not
    # need to do anything
    ctx.actions.write(ctx.outputs.executable, """
set -euox pipefail
find .
got=$(cat {json_file})
want='{{"data_files":["gazelle-runner.bash","gazelle"]}}'

if [ "$want" == "$got" ]; then
    echo 'PASS'
else
    exit 1
fi
""".format(
        json_file = ctx.outputs.json.short_path,
    ))

    runfiles = ctx.runfiles(files = [ctx.outputs.json], collect_data = True)

    return [DefaultInfo(
        files = depset([ctx.outputs.json, ctx.outputs.executable]),
        runfiles = runfiles,
    )]

proto_gazelle_data_test = rule(
    implementation = _proto_gazelle_data_test_impl,
    attrs = {
        "gazelle": attr.label(
            providers = [DefaultInfo],
            mandatory = True,
        ),
        "data": attr.label_list(
            allow_files = True,
        ),
    },
    outputs = {
        "json": "%{name}.json",
    },
    test = True,
)
