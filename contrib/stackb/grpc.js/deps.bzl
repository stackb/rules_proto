
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def grpc_js_deps():
    existing = native.existing_rules()  
    
    if "com_github_stackb_grpc_js" not in existing:
        http_archive(
            name = "com_github_stackb_grpc_js",
            urls = ["https://github.com/stackb/grpc.js/archive/c94ef115b4e8eea526d5b54b829cfc7542f39bc5.tar.gz"],
            strip_prefix = "grpc.js-c94ef115b4e8eea526d5b54b829cfc7542f39bc5",
            sha256 = "bf3b7fca7803a9187e6d6780089cad593997c46d76c5d78ba3202ce8b5e424b2",
        )
