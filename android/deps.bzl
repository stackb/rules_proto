load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def android_proto_compile_deps(
    rules_android_version = "60f03a20cefbe1e110ae0ac7f25359822e9ea24a",
    rules_android_sha256 = "4305b6cf6b098752a19fdb1abdc9ae2e069f5ff61359bfc3c752e4b4c862d18e",
    protobuf_lite_version = "5e8916e881c573c5d83980197a6f783c132d4276",
    protobuf_lite_sha256 = "d35902fb3cbe9afa67aad4e615a8224d0a531b8c06d32e100bdb235244748a3d",
    gmaven_tag = "20180927-1",
    gmaven_sha256 = "ddaa0f5811253e82f67ee637dc8caf3989e4517bac0368355215b0dcfa9844d6",
):
    existing = native.existing_rules()
    if "build_bazel_rules_android" not in existing:
        http_archive(
            name = "build_bazel_rules_android",
            urls = ["https://github.com/bazelbuild/rules_android/archive/%s.tar.gz" % rules_android_version],
            strip_prefix = "rules_android-%s" % rules_android_version,
            sha256 = rules_android_sha256,
        )
    if "com_google_protobuf_lite" not in existing:
        http_archive(
            name = "com_google_protobuf_lite",
            urls = ["https://github.com/protocolbuffers/protobuf/archive/%s.tar.gz" % protobuf_lite_version],
            strip_prefix = "protobuf-%s" % protobuf_lite_version,
            sha256 = protobuf_lite_sha256,
        )

    # if "android_sdk_tools" not in existing:
    #     http_archive(
    #         name = "android_sdk_tools",
    #         urls = ["https://dl.google.com/android/repository/sdk-tools-linux-4333796.zip"],
    #         #strip_prefix = "protobuf-%s" % protobuf_lite_version,
    #         #sha256 = protobuf_lite_sha256,
    #     )
    if "gmaven_rules" not in existing:
        http_archive(
            name = "gmaven_rules",
            strip_prefix = "gmaven_rules-%s" % gmaven_tag,
            urls = ["https://github.com/bazelbuild/gmaven_rules/archive/%s.tar.gz" % gmaven_tag],
            sha256 = gmaven_sha256,
        )

