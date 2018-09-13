workspace(name = "com_github_stackb_rules_grpc")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

http_archive(
    name = "com_github_grpc_grpc",
    strip_prefix = "grpc-1.15.0",
    url = "https://github.com/grpc/grpc/archive/v1.15.0.tar.gz",
    sha256 = "013cc34f3c51c0f87e059a12ea203087a7a15dca2e453295345e1d02e2b9634b",
)

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@//python:deps.bzl", "py_proto_deps")

py_proto_deps()

load("@//java:deps.bzl", "java_proto_deps")

java_proto_deps()

PROTOTOOL_VERSION="v1.2.0"

http_file(
    name = "com_github_uber_prototool_linux",
    urls = ["https://github.com/uber/prototool/releases/download/%s/prototool-Linux-x86_64" % PROTOTOOL_VERSION],
    sha256 = "cdbe781f8c3e3ed0a40490c33d0b8490fd5cb5a2c7912306f0d016878f6e26bd",
)

http_file(
    name = "com_github_uber_prototool_darwin",
    urls = ["https://github.com/uber/prototool/releases/download/%s/prototool-Darwin-x86_64" % PROTOTOOL_VERSION],
    sha256 = "cdbe781f8c3e3ed0a40490c33d0b8490fd5cb5a2c7912306f0d016878f6e26bd",
)
