"deps.bzl contains core repo dependencies."

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def io_bazel_rules_go():
    # Release: v0.28.0
    # TargetCommitish: release-0.28
    # Date: 2021-07-06 23:21:45 +0000 UTC
    # URL: https://github.com/bazelbuild/rules_go/releases/tag/v0.28.0
    # Size: 687092 (687 kB)
    maybe(
        http_archive,
        name = "io_bazel_rules_go",
        sha256 = "38171ce619b2695fa095427815d52c2a115c716b15f4cd0525a88c376113f584",
        strip_prefix = "rules_go-0.28.0",
        urls = ["https://github.com/bazelbuild/rules_go/archive/v0.28.0.tar.gz"],
    )

def bazel_gazelle():
    # Release: v0.22.3
    # Commit: release-0.22
    # Date: 2020-12-23 18:58:38 +0000 UTC
    # URL: https://github.com/bazelbuild/bazel-gazelle/releases/tag/v0.22.3
    # Branch: master
    # Commit: e4496b956eb46bdf8c9bf95b8d1d85e3a086c4be
    # Date: 2021-02-01 16:04:14 +0000 UTC
    # URL: https://github.com/bazelbuild/bazel-gazelle/commit/e4496b956eb46bdf8c9bf95b8d1d85e3a086c4be
    #
    # Upgrade Google go-cmp to v0.5.4 (#986)
    # Size: 1633162 (1.6 MB)
    maybe(
        http_archive,
        name = "bazel_gazelle",
        sha256 = "7359a9a7071dff5343a52626fce1aae2a78936d3004ef89d038daaefd3fd6608",
        strip_prefix = "bazel-gazelle-e4496b956eb46bdf8c9bf95b8d1d85e3a086c4be",
        urls = ["https://github.com/bazelbuild/bazel-gazelle/archive/e4496b956eb46bdf8c9bf95b8d1d85e3a086c4be.tar.gz"],
    )

def rules_proto():
    # Branch: master
    # Commit: f7a30f6f80006b591fa7c437fe5a951eb10bcbcf
    # Date: 2021-02-09 14:25:06 +0000 UTC
    # URL: https://github.com/bazelbuild/rules_proto/commit/f7a30f6f80006b591fa7c437fe5a951eb10bcbcf
    #
    # Merge pull request #77 from Yannic/proto_descriptor_set_rule
    #
    # Create proto_descriptor_set
    # Size: 14397 (14 kB)
    maybe(
        http_archive,
        name = "rules_proto",
        sha256 = "9fc210a34f0f9e7cc31598d109b5d069ef44911a82f507d5a88716db171615a8",
        strip_prefix = "rules_proto-f7a30f6f80006b591fa7c437fe5a951eb10bcbcf",
        urls = ["https://github.com/bazelbuild/rules_proto/archive/f7a30f6f80006b591fa7c437fe5a951eb10bcbcf.tar.gz"],
    )

def build_bazel_rules_swift():
    # Release: 0.18.0
    # Commit: master
    # Date: 2021-01-04 23:36:38 +0000 UTC
    # URL: https://github.com/bazelbuild/rules_swift/releases/tag/0.18.0
    # Branch: master
    # Commit: dadd12190182530cf6f91ca7f9e70391644ce502
    # Date: 2021-02-08 21:24:10 +0000 UTC
    # URL: https://github.com/bazelbuild/rules_swift/commit/dadd12190182530cf6f91ca7f9e70391644ce502
    #
    # Don't re-export the modules imported by a Swift generated header.
    #
    # This was an unintentional change in behavior from https://github.com/bazelbuild/rules_swift/commit/5f51ca9c5149122f41cada6122c61788d880fee9; this puts us back to the original behavior, but leaves an API in place for finer-grained control over re-exporting modules in the future. (But the BUILD rules today don't really have the flexibility to support it yet.)
    #
    # PiperOrigin-RevId: 356338982
    # (cherry picked from commit f45eea8c02a87c3077e5209f471fe4a193b5b0ba)
    # Size: 157518 (158 kB)
    maybe(
        http_archive,
        name = "build_bazel_rules_swift",
        sha256 = "1f5499bb053736cda8905d89aac42e98011bbe9ca93b774a40c04759f045d7bf",
        strip_prefix = "rules_swift-dadd12190182530cf6f91ca7f9e70391644ce502",
        urls = ["https://github.com/bazelbuild/rules_swift/archive/dadd12190182530cf6f91ca7f9e70391644ce502.tar.gz"],
    )

def io_bazel_rules_closure():
    # Commit: 2c59208867759800a37d0f008c3a4398af4c0cb2
    # Date: 2021-10-18 15:15:18 +0000 UTC
    # URL: https://github.com/bazelbuild/rules_closure/commit/2c59208867759800a37d0f008c3a4398af4c0cb2
    #
    # Upgrade com_google_template_soy_jssrc to a1c02e60ae88ed1b7db92722ea25ac7d396514fc
    # https://github.com/bazelbuild/rules_closure/pull/536
    # Size: 453381 (453 kB)
    maybe(
        http_archive,
        name = "io_bazel_rules_closure",
        sha256 = "50096b6be0052055ba4f0577d8aa3d82adf077377ffa86e2b7a67a335442f01b",
        strip_prefix = "rules_closure-2c59208867759800a37d0f008c3a4398af4c0cb2",
        urls = ["https://github.com/bazelbuild/rules_closure/archive/2c59208867759800a37d0f008c3a4398af4c0cb2.tar.gz"],
    )

def com_github_stackb_grpc_js():
    # Branch: master
    # Commit: beb6ac3b43247816c1a1ebf741ebf0c98203414a
    # Date: 2021-02-09 22:02:48 +0000 UTC
    # URL: https://github.com/stackb/grpc.js/commit/beb6ac3b43247816c1a1ebf741ebf0c98203414a
    #
    # Fix dangling build link (#7)
    #
    # * Fix dangling build ref
    # Size: 30483 (30 kB)
    maybe(
        http_archive,
        name = "com_github_stackb_grpc_js",
        sha256 = "f9cb4d932badc71d90a89263eabc93551923bb5c621e0940c7cfeaa79ef02596",
        strip_prefix = "grpc.js-beb6ac3b43247816c1a1ebf741ebf0c98203414a",
        urls = ["https://github.com/stackb/grpc.js/archive/beb6ac3b43247816c1a1ebf741ebf0c98203414a.tar.gz"],
    )

def io_grpc_grpc_java():
    # Release: v1.35.0
    # Commit: master
    # Date: 2021-01-12 23:05:49 +0000 UTC
    # URL: https://github.com/grpc/grpc-java/releases/tag/v1.35.0
    # Branch: master
    # Commit: 7f7821c616598ce4e33d2045c5641b2348728cb8
    # Date: 2021-02-10 00:56:26 +0000 UTC
    # URL: https://github.com/grpc/grpc-java/commit/7f7821c616598ce4e33d2045c5641b2348728cb8
    #
    # interop-testing: add fake altsHandshakerService for test (#7847)
    # Size: 2337953 (2.3 MB)
    maybe(
        http_archive,
        name = "io_grpc_grpc_java",
        sha256 = "82b3cf09f98a5932e1b55175aaec91b2a3f424eec811e47b2a3be533044d9afb",
        strip_prefix = "grpc-java-7f7821c616598ce4e33d2045c5641b2348728cb8",
        urls = ["https://github.com/grpc/grpc-java/archive/7f7821c616598ce4e33d2045c5641b2348728cb8.tar.gz"],
    )

def build_bazel_rules_nodejs():
    # This is a snapshot build from Greg Magolan that should support
    # "multilinker" (ability for nodejs_binary) to have deps from different
    # package.jsons
    maybe(
        http_archive,
        name = "build_bazel_rules_nodejs",
        sha256 = "8617ef45e5691e454835031541f404c84afbab0ad7f3ef62a853a45cd70b7df7",
        urls = [
            "https://github.com/aspect-dev/rules_nodejs-builds/raw/3.1.0+bd9eeb0e/build_bazel_rules_nodejs-snapshot_builds-snapshot.tar.gz",
        ],
    )

def rules_pkg():
    maybe(
        http_archive,
        name = "rules_pkg",
        urls = [
            "https://github.com/bazelbuild/rules_pkg/releases/download/0.2.6-1/rules_pkg-0.2.6.tar.gz",
            "https://mirror.bazel.build/github.com/bazelbuild/rules_pkg/releases/download/0.2.6/rules_pkg-0.2.6.tar.gz",
        ],
        sha256 = "aeca78988341a2ee1ba097641056d168320ecc51372ef7ff8e64b139516a4937",
    )

def rules_jvm_external():
    RULES_JVM_EXTERNAL_TAG = "4.0"
    RULES_JVM_EXTERNAL_SHA = "31701ad93dbfe544d597dbe62c9a1fdd76d81d8a9150c2bf1ecf928ecdf97169"
    maybe(
        http_archive,
        name = "rules_jvm_external",
        strip_prefix = "rules_jvm_external-%s" % RULES_JVM_EXTERNAL_TAG,
        sha256 = RULES_JVM_EXTERNAL_SHA,
        url = "https://github.com/bazelbuild/rules_jvm_external/archive/%s.zip" % RULES_JVM_EXTERNAL_TAG,
    )

def rules_codeowners():
    # Branch: master
    # Commit: 27fe3bbe6e5b0df196e360fc9e081835f22a10be
    # Date: 2020-06-10 13:20:02 +0000 UTC
    # URL: https://github.com/zegl/rules_codeowners/commit/27fe3bbe6e5b0df196e360fc9e081835f22a10be
    #
    # Allow user to override the comment in the generated file
    #
    # Useful if they have some specific instructions to put there about how to update it
    # Size: 7198 (7.2 kB)
    maybe(
        http_archive,
        name = "rules_codeowners",
        sha256 = "2972e61f08dd41bb00fa4c05a36949a51b32921ae714aa900d9d755ad00533f5",
        strip_prefix = "rules_codeowners-27fe3bbe6e5b0df196e360fc9e081835f22a10be",
        urls = ["https://github.com/zegl/rules_codeowners/archive/27fe3bbe6e5b0df196e360fc9e081835f22a10be.tar.gz"],
    )

def com_github_grpc_grpc():
    # Release: v1.35.0
    # Commit: v1.35.x
    # Date: 2021-01-19 18:07:57 +0000 UTC
    # URL: https://github.com/grpc/grpc/releases/tag/v1.35.0
    # Branch: master
    # Commit: 5f759fcd1f602b38004b948b071f8b5726a9a4b1
    # Date: 2021-02-09 05:46:27 +0000 UTC
    # URL: https://github.com/grpc/grpc/commit/5f759fcd1f602b38004b948b071f8b5726a9a4b1
    #
    # Merge pull request #25384 from gnossen/fix_interop_typo
    #
    # Fix Interop Client Typo
    # Size: 7899154 (7.9 MB)
    maybe(
        http_archive,
        name = "com_github_grpc_grpc",
        sha256 = "e6c6b1ac9ba2257c93e49c98ef4fc96b2e2a1cdd90782a919f60e23fa8c2428b",
        strip_prefix = "grpc-5f759fcd1f602b38004b948b071f8b5726a9a4b1",
        urls = ["https://github.com/grpc/grpc/archive/5f759fcd1f602b38004b948b071f8b5726a9a4b1.tar.gz"],
    )

def rules_cc():
    # Branch: master
    # Commit: 40548a2974f1aea06215272d9c2b47a14a24e556
    # Date: 2021-02-05 12:29:43 +0000 UTC
    # URL: https://github.com/bazelbuild/rules_cc/commit/40548a2974f1aea06215272d9c2b47a14a24e556
    #
    # Automatic code cleanup.
    #
    # PiperOrigin-RevId: 355825197
    # Change-Id: I8acfc20228816c04fcf48bfcc435cbde2b1fb608
    # Size: 129521 (130 kB)
    maybe(
        http_archive,
        name = "rules_cc",
        sha256 = "cb8ce8a25464b2a8536450971ad1b45ee309491c1f5e052a611b9e249cfdd35d",
        strip_prefix = "rules_cc-40548a2974f1aea06215272d9c2b47a14a24e556",
        urls = ["https://github.com/bazelbuild/rules_cc/archive/40548a2974f1aea06215272d9c2b47a14a24e556.tar.gz"],
    )
