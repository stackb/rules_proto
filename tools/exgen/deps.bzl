load("@bazel_gazelle//:deps.bzl", "go_repository")

def exgen_deps(**kwargs):
    go_repository(
        name = "com_github_urfave_cli",
        importpath = "github.com/urfave/cli",
        urls = ["https://github.com/urfave/cli/archive/cfb38830724cc34fedffe9a2a29fb54fa9169cd1.tar.gz"],
        strip_prefix = "cli-cfb38830724cc34fedffe9a2a29fb54fa9169cd1/",
        build_file_generation = "on",
        build_file_proto_mode = "disable",
    )
