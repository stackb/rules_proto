"""proto_gazelle_data_test is a regression test.

See https://github.com/stackb/rules_proto/issues/342.

"""

# This test script asserts that all expected files are present in runfiles.
# Specifically, as a regression test that items listed in gazelle.data are
# present (./genfile_should_be_present_in_gazelle_data_runfiles.txt).
SCRIPT = """
set -x

read -r -d '' WANT <<'EOF'
.
./config.yaml
./external
./external/go_sdk
./external/go_sdk/bin
./external/go_sdk/bin/go
./external/googleapis
./external/googleapis/imports.csv
./gazelle
./gazelle-protobuf_
./gazelle-protobuf_/gazelle-protobuf
./gazelle-runner.bash
./genfile_should_be_present_in_gazelle_data_runfiles.txt
./proto_repository_data_test
EOF

GOT=$(find . | sort)

if [ "$WANT" == "$GOT" ]; then
    echo 'PASS'
else
    echo "WANT:\n$WANT"
    echo "GOT:\n$GOT"
    exit 1
fi
"""

def _proto_gazelle_data_test_impl(ctx):
    ctx.actions.write(ctx.outputs.executable, SCRIPT)

    runfiles = ctx.runfiles().merge(ctx.attr.gazelle[DefaultInfo].default_runfiles)

    return [DefaultInfo(
        files = depset([ctx.outputs.executable]),
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
    test = True,
)
