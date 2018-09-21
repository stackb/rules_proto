
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def grpc_web_deps():
    existing = native.existing_rules()  
    
    if "com_github_grpc_grpc_web" not in existing:
        http_archive(
            name = "com_github_grpc_grpc_web",
            urls = ["https://github.com/grpc/grpc-web/archive/92aa9f8fc8e7af4aadede52ea075dd5790a63b62.tar.gz"],
            strip_prefix = "grpc-web-92aa9f8fc8e7af4aadede52ea075dd5790a63b62",
            sha256 = "f4996205e6d1d72e2be46f1bda4d26f8586998ed42021161322d490537d8c9b9",
        )
