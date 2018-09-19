workspace(name = "com_github_stackb_rules_grpc")

# =========================================

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

# =========================================

http_archive(
    name = "com_github_grpc_grpc",
    strip_prefix = "grpc-1.15.0",
    url = "https://github.com/grpc/grpc/archive/v1.15.0.tar.gz",
    sha256 = "013cc34f3c51c0f87e059a12ea203087a7a15dca2e453295345e1d02e2b9634b",
)

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

# =========================================

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

# ===========

load("@//closure:deps.bzl", "closure_proto_library_deps")

closure_proto_library_deps()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories(
    omit_com_google_protobuf = True,
    #omit_com_google_code_findbugs_jsr305 = True,
    # omit_com_google_errorprone_error_prone_annotations = True,
)

# =========================================

load("@//go:deps.bzl", "go_proto_deps")

go_proto_deps()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

# =========================================

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")
grpc_java_repositories(
    omit_com_google_protobuf = True,
)

# =========================================

local_repository(
    name = "io_bazel_rules_dart",
    path = "/home/pcj/github/dart-lang/rules_dart",
)

load("@//dart:deps.bzl", "dart_proto_deps", "dart_pub_deps")

dart_proto_deps()

dart_pub_deps(
    name = "dart_pub_deps_protoc_plugin",
    spec = "//dart:pubspec.yaml",
    override = {
        "path": "1.6.2",
        "analyzer": "0.32.5",
        "crypto": "2.0.6",
        "async": "2.0.8",
        "fixnum": "0.10.8",
        "collection": "1.14.11",
        "dart_style": "1.1.3",
        "source_span": "1.4.1",
        "args": "1.5.0",
    },
    verbose = 0,
)

load("@dart_pub_deps_protoc_plugin//:deps.bzl", "pub_deps")
pub_deps(verbose = 0)

load("@io_bazel_rules_dart//dart/build_rules:repositories.bzl", "dart_repositories")
dart_repositories()

load("@io_bazel_rules_dart//dart/build_rules/internal:pub.bzl", "pub_repository")

pub_repository(
        name = "vendor_isolate",
        output = ".",
        package = "isolate",
        version = "2.0.2",
    )

