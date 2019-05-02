load(
    "//:deps.bzl",
    "get_ref",
    "get_sha256",
)
load(
    "//protobuf:deps.bzl",
    "protobuf",
)

# Special thing to get around maven jar issues
load("//closure:buildozer_http_archive.bzl", "buildozer_http_archive")

def io_bazel_rules_closure(**kwargs):
    name = "io_bazel_rules_closure"

    # PR#361 - includes closure_js_library.library_level_checks
    ref = get_ref(name, "e86d8021f22277fe129a572cd019e846243d4531", kwargs)  # PR #361
    sha256 = get_sha256(name, "481b6b522c2894906380b4b9c008b2c37ab86eeb182229d75bf453db89ed79bc", kwargs)

    # ref = get_ref(name, "50d3dc9e6d27a5577a0f95708466718825d579f4", kwargs) # HEAD April 2019
    # sha256 = get_sha256(name, "1c05fea22c9630cf1047f25d008780756373a60ddd4d2a6993cf9858279c5da6", kwargs)

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
    protobuf(**kwargs)

def closure_proto_library(**kwargs):
    closure_proto_compile(**kwargs)
    io_bazel_rules_closure(**kwargs)
