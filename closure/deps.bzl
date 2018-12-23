load(
    "//:deps.bzl",
    "com_google_protobuf",
    "get_ref",
    "get_sha256",
)

# Special thing to get around maven jar issues
load("//closure:buildozer_http_archive.bzl", "buildozer_http_archive")

def io_bazel_rules_closure(**kwargs):
    name = "io_bazel_rules_closure"
    ref = get_ref(name, "9889e2348259a5aad7e805547c1a0cf311cfcd91", kwargs)
    sha256 = get_sha256(name, "43c9b882fa921923bcba764453f4058d102bece35a37c9f6383c713004aacff1", kwargs)

    if "io_bazel_rules_closure" not in native.existing_rules():
        buildozer_http_archive(
            name = "io_bazel_rules_closure",
            urls = ["https://github.com/bazelbuild/rules_closure/archive/%s.tar.gz" % ref],
            sha256 = sha256,
            strip_prefix = "rules_closure-" + ref,
            label_list = ["//...:%java_binary", "//...:%java_library"],
            replace_deps = {
                "@com_google_code_findbugs_jsr305": "@com_google_code_findbugs_jsr305_3_0_0",
                "@com_google_errorprone_error_prone_annotations": "@com_google_errorprone_error_prone_annotations_2_1_3",
            },
            sed_replacements = {
                "closure/repositories.bzl": [
                    "s|com_google_code_findbugs_jsr305|com_google_code_findbugs_jsr305_3_0_0|g",
                    "s|com_google_errorprone_error_prone_annotations|com_google_errorprone_error_prone_annotations_2_1_3|g",
                ],
            },
        )

def closure_proto_compile(**kwargs):
    com_google_protobuf(**kwargs)

def closure_proto_library(**kwargs):
    closure_proto_compile(**kwargs)
    io_bazel_rules_closure(**kwargs)
