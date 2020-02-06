# Copied from build_bazel_rules_nodejs and changed the strip_prefix

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def npm_bazel_typescript():
    if "npm_bazel_typescript" not in native.existing_rules():
        
        # Branch: master
        # Commit: 7f7a45b603f8839e3781d993c24ed6f11050f6bf
        # Date: 2020-02-05 20:44:10 +0000 UTC
        # URL: https://github.com/bazelbuild/rules_nodejs/commit/7f7a45b603f8839e3781d993c24ed6f11050f6bf
        # 
        # fix: html script injection is broken on windows
        # 
        # `html-insert-assets` version `0.5.0` in includes the following fix: https://github.com/jbedard/html-insert-assets/pull/6 which is needed for better Windows support.
        # 
        # Closes #1604
        # Size: 4636042 (4.6 MB)
        http_archive(
            name = "npm_bazel_typescript",
            sha256 = "c237630c9ae5122e46cd1e11ecf3c04316978058dd30825330784a6f464c75c4",
            strip_prefix = "rules_nodejs-7f7a45b603f8839e3781d993c24ed6f11050f6bf/packages/typescript/src",
            urls = ["https://github.com/bazelbuild/rules_nodejs/archive/7f7a45b603f8839e3781d993c24ed6f11050f6bf.tar.gz"],
        )