

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def go_proto_compile_deps():
    existing = native.existing_rules()

    if "io_bazel_rules_go" not in existing:
        http_archive(
            name = "io_bazel_rules_go",
            urls = ["https://github.com/bazelbuild/rules_go/releases/download/0.15.3/rules_go-0.15.3.tar.gz"],
            sha256 = "97cf62bdef33519412167fd1e4b0810a318a7c234f5f8dc4f53e2da86241c492",
        )

def go_proto_library_deps():
    go_proto_compile_deps()