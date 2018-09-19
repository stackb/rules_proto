load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")
load("@//closure:buildozer_http_archive.bzl", "buildozer_http_archive")

RULES_CLOSURE_VERSION = "1e12aa5612d758daf2df339991c8d187223a7ee6"
RULES_CLOSURE_SHA256 = "663424d34fd067a8d066308eb2887fcaba36d73b354669ec1467498726a6b82c"

def closure_proto_compile_deps():
    pass

def closure_proto_library_deps():
    closure_proto_compile_deps()

    existing = native.existing_rules()

    if "com_github_bazelbuild_buildtools_buildozer_linux" not in existing:
        http_file(
            name = "com_github_bazelbuild_buildtools_buildozer_linux",
            urls = ["https://github.com/bazelbuild/buildtools/releases/download/0.15.0/buildozer"],
            sha256 = "be07a37307759c68696c989058b3446390dd6e8aa6fdca6f44f04ae3c37212c5",
        )      

    if "com_github_bazelbuild_buildtools_buildozer_darwin" not in existing:
        http_file(
            name = "com_github_bazelbuild_buildtools_buildozer_darwin",
            urls = ["https://github.com/bazelbuild/buildtools/releases/download/0.15.0/buildozer.osx"],
            sha256 = "294357ff92e7bb36c62f964ecb90e935312671f5a41a7a9f2d77d8d0d4bd217d",
        )      

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
