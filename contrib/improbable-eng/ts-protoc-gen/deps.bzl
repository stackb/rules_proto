
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def ts_proto_deps():
    existing = native.existing_rules()  
    
    if "build_bazel_rules_nodejs" not in existing:
        http_archive(
            name = "build_bazel_rules_nodejs",
            urls = ["https://github.com/bazelbuild/rules_nodejs/archive/d334fd8e2274fb939cf447106dced97472534e80.tar.gz"],
            strip_prefix = "rules_nodejs-d334fd8e2274fb939cf447106dced97472534e80",
            sha256 = "5c69bae6545c5c335c834d4a7d04b888607993027513282a5139dbbea7166571",
        )

    if "build_bazel_rules_typescript" not in existing:
        http_archive(
            name = "build_bazel_rules_typescript",
            urls = ["https://github.com/bazelbuild/rules_typescript/archive/3488d4fb89c6a02d79875d217d1029182fbcd797.tar.gz"],
            strip_prefix = "rules_typescript-3488d4fb89c6a02d79875d217d1029182fbcd797",
            sha256 = "22ebe19999ce34de2f0329d29c7cac1cccd449cd61d0813aa0e633ac8dfaef80",
        )

    if "io_bazel_rules_webtesting" not in existing:
        http_archive(
            name = "io_bazel_rules_webtesting",
            urls = ["https://github.com/bazelbuild/rules_webtesting/archive/e417b122a3d1e8f8a4cc09b1b05e2a5f52c8ecbb.tar.gz"],
            strip_prefix = "rules_webtesting-e417b122a3d1e8f8a4cc09b1b05e2a5f52c8ecbb",
            sha256 = "76af36aac2aed3ca6984e0c45d1204582243fa6b81165bf0e42eacfffed9c769",
        )

    if "ts_protoc_gen" not in existing:
        http_archive(
            name = "ts_protoc_gen",
            urls = ["https://github.com/improbable-eng/ts-protoc-gen/archive/67e0c93fa4e29539ec57c97d23116f366ffb94e7.tar.gz"],
            strip_prefix = "ts-protoc-gen-67e0c93fa4e29539ec57c97d23116f366ffb94e7",
            sha256 = "03b87365116b3e829b648c57bf17b8d252b682c33264de5663e7f410c9ad0e55",
        )

    if "bazel_gazelle" not in existing:
        http_archive(
            name = "bazel_gazelle",
            urls = ["https://github.com/bazelbuild/bazel-gazelle/releases/download/0.14.0/bazel-gazelle-0.14.0.tar.gz"],
            sha256 = "c0a5739d12c6d05b6c1ad56f2200cb0b57c5a70e03ebd2f7b87ce88cabf09c7b",
        )
