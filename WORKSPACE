workspace(name = "com_github_stackb_rules_grpc")

# =========================================

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

# =========================================
load("//:deps.bzl", "grpc_deps")

grpc_deps()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", com_github_grpc_grpc_bazel_grpc_deps = "grpc_deps")

com_github_grpc_grpc_bazel_grpc_deps()

# =========================================

load("//python:deps.bzl", "py_proto_compile_deps")

py_proto_compile_deps()

load("@io_bazel_rules_python//python:pip.bzl", "pip_repositories", "pip_import")

pip_repositories()

pip_import(
   name = "grpc_py_deps",
   requirements = "//python:requirements.txt",
)

load("@grpc_py_deps//:requirements.bzl", "pip_install")
pip_install()

# =========================================

load("//java:deps.bzl", "java_grpc_library_deps")

java_grpc_library_deps()

# ===========

load("//closure:deps.bzl", "closure_proto_library_deps")

closure_proto_library_deps()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories(
    omit_com_google_protobuf = True,
)

# =========================================

load("//node:deps.bzl", "node_proto_library_deps")

node_proto_library_deps()

load("@org_pubref_rules_node//node:rules.bzl", "node_repositories")

node_repositories()

load("@org_pubref_rules_node//node:rules.bzl", "yarn_modules")

yarn_modules(
    name = "proto_node_modules",
    deps = {
        "google-protobuf": "3.6.1",
    },
)

yarn_modules(
    name = "grpc_node_modules",
    deps = {
        "google-protobuf": "3.6.1",
        "grpc": "1.15.1",
    },
)

# =========================================

load("//go:deps.bzl", "go_proto_library_deps")

go_proto_library_deps()

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

load("//dart:deps.bzl", "dart_proto_deps", "dart_pub_deps")

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

# =========================================

load("//github.com/stackb/grpc.js:deps.bzl", "grpc_js_deps")

grpc_js_deps()

# =========================================

load("//github.com/grpc/grpc-web:deps.bzl", "grpc_web_deps")

grpc_web_deps()

# =========================================

load("//github.com/improbable-eng/ts-protoc-gen:deps.bzl", "ts_proto_deps")

ts_proto_deps()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

load("@build_bazel_rules_nodejs//:defs.bzl", "node_repositories")

node_repositories(
    package_json = ["@ts_protoc_gen//:package.json"],
)

load("@build_bazel_rules_typescript//:defs.bzl", "ts_setup_workspace")

ts_setup_workspace()

load("@io_bazel_rules_webtesting//web:repositories.bzl", "browser_repositories", "web_test_repositories")

web_test_repositories()

load("@build_bazel_rules_nodejs//:defs.bzl", "npm_install")

npm_install(
    name = "deps",
    package_json = "@ts_protoc_gen//:package.json",
    package_lock_json = "@ts_protoc_gen//:package-lock.json",
)

# =======================================

load("//ruby:deps.bzl", "ruby_proto_library_deps")

ruby_proto_library_deps()

load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_register_toolchains")

ruby_register_toolchains()

load("@com_github_yugui_rules_ruby//ruby/private:bundle.bzl", "bundle_install")
bundle_install(
    name = "routeguide_gems_bundle",
    gemfile = "//ruby:Gemfile",
    gemfile_lock = "//ruby:Gemfile.lock",
)
