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

load("@io_bazel_rules_python//python:pip.bzl", "pip_repositories", "pip_import")

pip_repositories()

pip_import(
   name = "grpc_py_deps",
   requirements = "//python:requirements.txt",
)

load("@grpc_py_deps//:requirements.bzl", "pip_install")
pip_install()

# =========================================

load("@//java:deps.bzl", "java_proto_deps")

java_proto_deps()

load("@//go:deps.bzl", "go_proto_deps")

go_proto_deps()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")
grpc_java_repositories(
    omit_com_google_protobuf = True,
)