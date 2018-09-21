
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def grpc_web_library_deps():
    existing = native.existing_rules()  
    
    if "com_github_improbable_eng_grpc_web" not in existing:
        http_archive(
            name = "com_github_improbable_eng_grpc_web",
            urls = ["https://github.com/improbable-eng/grpc-web/archive/b89ec6300fb9ce3f604ff02c22e1c106fe29a3b9.tar.gz"],
            strip_prefix = "grpc-web-b89ec6300fb9ce3f604ff02c22e1c106fe29a3b9",
            sha256 = "5d69bae6545c5c335c834d4a7d04b888607993027513282a5139dbbea7166571",
        )
