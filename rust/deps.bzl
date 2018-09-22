load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

libc_BUILD = """
"""

def rust_proto_library_deps(
  rules_rust_version = "88022d175adb48aa5f8904f95dfc716c543b3f1e",
  rules_rust_sha256 = "d9832945f0fa7097ee548bd6fecfc814bd19759561dd7b06723e1c6a1879aa71", 
):
    existing = native.existing_rules()

    if "io_bazel_rules_rust" not in existing:
        http_archive(
            name = "io_bazel_rules_rust",
            urls = ["https://github.com/bazelbuild/rules_rust/archive/%s.tar.gz" % rules_rust_version],
            sha256 = rules_rust_sha256,
            strip_prefix = "rules_rust-" + rules_rust_version,
        )
        # native.local_repository(
        #     name = "io_bazel_rules_rust",
        #     path = "/home/pcj/github/bazelbuild/rules_rust",
        # )

    # if "libc" not in existing:
    #     http_archive(
    #         name = "libc",
    #         urls = ["https://github.com/rust-lang/libc/archive/%s.tar.gz" % libc_version],
    #         sha256 = libc_sha256,
    #         #tag = "0.2.20",
    #         #build_file_content = libc_BUILD,
    #     )
