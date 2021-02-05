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
        build_file = "@build_stack_rules_proto//third_party:BUILD.bazel.zlib",
    )
