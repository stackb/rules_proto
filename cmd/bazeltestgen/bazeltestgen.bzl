load("@io_bazel_rules_go//go/tools/bazel_testing:def.bzl", "go_bazel_test")

def _bazeltestgen_impl(ctx):
    config_json = ctx.outputs.json
    output_bazel_test = ctx.outputs.bazel_test

    config = struct(
        out = output_bazel_test.path,
        files = [f.path for f in ctx.files.srcs],
    )

    ctx.actions.write(
        output = config_json,
        content = config.to_json(),
    )

    ctx.actions.run(
        mnemonic = "BazelTestGenerate",
        progress_message = "Generating %s test" % ctx.attr.name,
        executable = ctx.file._bazeltestgen,
        arguments = ["--config_json=%s" % config_json.path],
        inputs = [config_json] + ctx.files.srcs,
        outputs = [output_bazel_test],
    )

    return [DefaultInfo(
        files = depset([config_json, output_bazel_test]),
    )]

bazeltestgen = rule(
    implementation = _bazeltestgen_impl,
    attrs = {
        "srcs": attr.label_list(
            doc = "Sources for the test txtar file",
            allow_files = True,
        ),
        "_bazeltestgen": attr.label(
            doc = "The bazeltestgen generator tool",
            default = "//cmd/bazeltestgen",
            allow_single_file = True,
            executable = True,
            cfg = "host",
        ),
    },
    outputs = {
        "json": "%{name}.json",
        "bazel_test": "%{name}_bazel_test.go",
    },
)

def gazelle_testdata_go_bazel_test(**kwargs):
    name = kwargs.pop("name")
    srcs = kwargs.pop("srcs", [])
    rule_files = kwargs.pop("rule_files", ["//:all_files"])

    bazeltestgen(
        name = name + "gen",
        srcs = srcs,
    )

    go_bazel_test(
        name = name,
        srcs = [name + "gen_bazel_test.go"],
        rule_files = rule_files,
        **kwargs
    )
