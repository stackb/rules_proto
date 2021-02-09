load("@io_bazel_rules_go//go/tools/bazel_testing:def.bzl", "go_bazel_test")

def proto_rule_bazel_test(**kwargs):
    name = kwargs.pop("name")
    deps = kwargs.pop("deps", [])

    go_bazel_test(
        name = name,
        srcs = [dep_name + "_bazel_test.go" for dep_name in deps],
        rule_files = kwargs.get("rule_files", ["//:all_files"]),
        **kwargs
    )
