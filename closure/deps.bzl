load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")
load("//closure:buildozer_http_archive.bzl", "buildozer_http_archive")

RULES_CLOSURE_VERSION = "1e12aa5612d758daf2df339991c8d187223a7ee6"
RULES_CLOSURE_SHA256 = "663424d34fd067a8d066308eb2887fcaba36d73b354669ec1467498726a6b82c"

def closure_proto_library_deps():
    existing = native.existing_rules()

    if "io_bazel_rules_closure" not in existing:
        buildozer_http_archive(
            name = "io_bazel_rules_closure",
            urls = ["https://github.com/bazelbuild/rules_closure/archive/%s.tar.gz" % RULES_CLOSURE_VERSION],
            sha256 = RULES_CLOSURE_SHA256,
            strip_prefix = "rules_closure-" + RULES_CLOSURE_VERSION,
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
