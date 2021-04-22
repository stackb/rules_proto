load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def com_google_protobuf():
    # Release: v3.14.0
    # TargetCommitish: 3.14.x
    # Date: 2020-11-13 20:53:39 +0000 UTC
    # URL: https://github.com/protocolbuffers/protobuf/releases/tag/v3.14.0
    # Size: 5319779 (5.3 MB)
    maybe(
        http_archive,
        name = "com_google_protobuf",
        sha256 = "d0f5f605d0d656007ce6c8b5a82df3037e1d8fe8b121ed42e536f569dec16113",
        strip_prefix = "protobuf-3.14.0",
        urls = ["https://github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz"],
    )

def io_bazel_rules_go():
    # Release: v0.24.11
    # TargetCommitish: release-0.24
    # Date: 2021-01-19 23:11:54 +0000 UTC
    # URL: https://github.com/bazelbuild/rules_go/releases/tag/v0.24.11
    # Size: 523890 (524 kB)
    maybe(
        http_archive,
        name = "io_bazel_rules_go",
        sha256 = "d2b5751d8ae55ac011540453cf9da49ee12b832d0a98ca8ffae99285abb481f7",
        strip_prefix = "rules_go-0.24.11",
        urls = ["https://github.com/bazelbuild/rules_go/archive/v0.24.11.tar.gz"],
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

def bazel_skylib():
    # Branch: master
    # Commit: f80bc733d4b9f83d427ce3442be2e07427b2cc8d
    # Date: 2021-01-29 18:38:17 +0000 UTC
    # URL: https://github.com/bazelbuild/bazel-skylib/commit/f80bc733d4b9f83d427ce3442be2e07427b2cc8d
    #
    # update owners (#289)
    # Size: 89591 (90 kB)
    http_archive(
        name = "bazel_skylib",
        sha256 = "ebdf850bfef28d923a2cc67ddca86355a449b5e4f38b0a70e584dc24e5984aa6",
        strip_prefix = "bazel-skylib-f80bc733d4b9f83d427ce3442be2e07427b2cc8d",
        urls = ["https://github.com/bazelbuild/bazel-skylib/archive/f80bc733d4b9f83d427ce3442be2e07427b2cc8d.tar.gz"],
    )

def rules_python():
    # Branch: master
    # Commit: c7e068d38e2fec1d899e1c150e372f205c220e27
    # Date: 2021-02-02 22:16:45 +0000 UTC
    # URL: https://github.com/bazelbuild/rules_python/commit/c7e068d38e2fec1d899e1c150e372f205c220e27
    #
    # pip: 20.3.3 -> 20.3.4 (#405)
    # Size: 2563550 (2.6 MB)
    maybe(
        http_archive,
        name = "rules_python",
        sha256 = "8cc0ad31c8fc699a49ad31628273529ef8929ded0a0859a3d841ce711a9a90d5",
        strip_prefix = "rules_python-c7e068d38e2fec1d899e1c150e372f205c220e27",
        urls = ["https://github.com/bazelbuild/rules_python/archive/c7e068d38e2fec1d899e1c150e372f205c220e27.tar.gz"],
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

    # Branch: master
    # Commit: a0761ed101b939e19d83b2da5f59034bffc19c12
    # Date: 2021-01-26 15:30:54 +0000 UTC
    # URL: https://github.com/bazelbuild/rules_proto/commit/a0761ed101b939e19d83b2da5f59034bffc19c12
    #
    # Merge pull request #81 from Yannic/patch-3
    #
    # Bump bazel-toolchains to 3.7.2
    # Size: 11622 (12 kB)
    maybe(
        http_archive,
        name = "rules_proto",
        sha256 = "2a20fd8af3cad3fbab9fd3aec4a137621e0c31f858af213a7ae0f997723fc4a9",
        strip_prefix = "rules_proto-a0761ed101b939e19d83b2da5f59034bffc19c12",
        urls = ["https://github.com/bazelbuild/rules_proto/archive/a0761ed101b939e19d83b2da5f59034bffc19c12.tar.gz"],
    )

def zlib():
    maybe(
        http_archive,
        name = "zlib",
        urls = [
            "https://mirror.bazel.build/zlib.net/zlib-1.2.11.tar.gz",
            "https://zlib.net/zlib-1.2.11.tar.gz",
        ],
        sha256 = "c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1",
        strip_prefix = "zlib-1.2.11",
        build_file = "@build_stack_rules_proto//third_party:zlib.BUILD",
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
    # Branch: master
    # Commit: 4c99be33856ce1b7b80f55a0e9a8345f559b6ef3
    # Date: 2021-01-29 00:11:54 +0000 UTC
    # URL: https://github.com/bazelbuild/rules_closure/commit/4c99be33856ce1b7b80f55a0e9a8345f559b6ef3
    #
    # CheckStrictDeps: resolve module paths against js module roots (#510)
    #
    # Added "es6_relative_imports_gen_srcs_bin" test case
    #
    # Tested against existing cases and compiling an external
    # codebase with the change applied
    # Size: 468337 (468 kB)
    maybe(
        http_archive,
        name = "io_bazel_rules_closure",
        sha256 = "4c98a6b8d2d81210f3e291b1c7c5034ab2e22e7870ab3e9603599c79833f7da3",
        strip_prefix = "rules_closure-4c99be33856ce1b7b80f55a0e9a8345f559b6ef3",
        urls = ["https://github.com/bazelbuild/rules_closure/archive/4c99be33856ce1b7b80f55a0e9a8345f559b6ef3.tar.gz"],
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