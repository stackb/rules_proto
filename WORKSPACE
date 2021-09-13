workspace(name = "build_stack_rules_proto")

# gazelle:repo bazel_gazelle

# ==================================================
# Foundational
# ==================================================
#
load(
    "//:deps.bzl",
    "bazel_gazelle",
    "bazel_skylib",
    "build_bazel_rules_swift",
    "com_github_grpc_grpc",
    "com_google_protobuf",
    "io_bazel_rules_go",
    "io_grpc_grpc_java",
    "rules_jvm_external",
    "rules_proto",
    "rules_python",
    "zlib",
)

io_bazel_rules_go()

bazel_gazelle()

bazel_skylib()

com_github_grpc_grpc()

build_bazel_rules_swift()

io_grpc_grpc_java()

rules_jvm_external()

# ==================================================
# Go
# ==================================================
#
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.16.2")

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

# ==================================================
# Go Tool Deps
# ==================================================

load("//:go_deps.bzl", "go_deps")

# gazelle:repository_macro go_deps.bzl%go_deps
go_deps()

# # ==================================================
# # Protobuf Core
# # ==================================================
# #
rules_proto()

rules_python()

com_google_protobuf()

zlib()

# # load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")

# # rules_proto_dependencies()

# # rules_proto_toolchains()

# #
# # Toolchains
# #
register_toolchains("//toolchain")

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

### Java

load("@rules_jvm_external//:defs.bzl", "maven_install")
load("@io_grpc_grpc_java//:repositories.bzl", "IO_GRPC_GRPC_JAVA_ARTIFACTS", "IO_GRPC_GRPC_JAVA_OVERRIDE_TARGETS", "grpc_java_repositories")

maven_install(
    artifacts = IO_GRPC_GRPC_JAVA_ARTIFACTS,
    generate_compat_repositories = True,
    override_targets = IO_GRPC_GRPC_JAVA_OVERRIDE_TARGETS,
    repositories = [
        "https://repo.maven.apache.org/maven2/",
    ],
)

load("@maven//:compat.bzl", "compat_repositories")

compat_repositories()

grpc_java_repositories()

# load("//rules:py_grpc_compile_deps.bzl", "py_grpc_compile_deps")

# py_grpc_compile_deps()

# # gazelle:repository_macro tools/laze/deps.bzl%laze_deps
# laze_deps()

# build_bazel_rules_swift()

# load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

# grpc_deps()

# io_bazel_rules_closure()

# load("@io_bazel_rules_closure//closure:repositories.bzl", "rules_closure_dependencies", "rules_closure_toolchains")

# rules_closure_dependencies(
#     omit_bazel_skylib = True,
#     omit_com_google_protobuf = True,
#     omit_zlib = True,
# )

# rules_closure_toolchains()

# com_github_stackb_grpc_js()

# io_grpc_grpc_java()

# load("@io_grpc_grpc_java//:repositories.bzl", "IO_GRPC_GRPC_JAVA_ARTIFACTS", "IO_GRPC_GRPC_JAVA_OVERRIDE_TARGETS", "grpc_java_repositories")

# grpc_java_repositories()

# # rules_pkg()

# # rules_codeowners()

# # =================================================================================================
# # java
# # =================================================================================================

# load("@build_stack_rules_proto//rules:java_grpc_library_deps.bzl", "java_grpc_library_deps")

# java_grpc_library_deps()

# load("@rules_jvm_external//:defs.bzl", "maven_install")

# maven_install(
#     artifacts = IO_GRPC_GRPC_JAVA_ARTIFACTS,
#     generate_compat_repositories = True,
#     override_targets = IO_GRPC_GRPC_JAVA_OVERRIDE_TARGETS,
#     repositories = [
#         "https://repo.maven.apache.org/maven2/",
#     ],
# )

# load("@maven//:compat.bzl", "compat_repositories")

# compat_repositories()

# grpc_java_repositories()

# # =================================================================================================
# # nodejs
# # =================================================================================================

# build_bazel_rules_nodejs()

# load("@build_bazel_rules_nodejs//:index.bzl", "yarn_install")

# yarn_install(
#     name = "google_protobuf_node_modules",
#     package_json = "@build_stack_rules_proto//plugins/nodejs/modules/google-protobuf:package.json",
#     package_path = "plugins/nodejs/modules/google-protobuf",
#     symlink_node_modules = False,
#     yarn_lock = "@build_stack_rules_proto//plugins/nodejs/modules/google-protobuf:yarn.lock",
# )

# yarn_install(
#     name = "grpc_js_node_modules",
#     package_json = "@build_stack_rules_proto//plugins/nodejs/modules/grpc-js:package.json",
#     package_path = "plugins/nodejs/modules/grpc-js",
#     symlink_node_modules = False,
#     yarn_lock = "@build_stack_rules_proto//plugins/nodejs/modules/grpc-js:yarn.lock",
# )

# yarn_install(
#     name = "grpc_tools_node_modules",
#     package_json = "@build_stack_rules_proto//plugins/nodejs/modules/grpc-tools:package.json",
#     package_path = "plugins/nodejs/modules/grpc-tools",
#     symlink_node_modules = False,
#     yarn_lock = "@build_stack_rules_proto//plugins/nodejs/modules/grpc-tools:yarn.lock",
# )

# yarn_install(
#     name = "example_routeguide_node_modules",
#     package_json = "@build_stack_rules_proto//example/routeguide/nodejs/modules:package.json",
#     symlink_node_modules = False,
#     yarn_lock = "@build_stack_rules_proto//example/routeguide/nodejs/modules:yarn.lock",
# )

# #
# # Core
# #
# load("//:repositories.bzl", "rules_proto_grpc_repos")
# rules_proto_grpc_repos()

# load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")
# rules_proto_dependencies()
# rules_proto_toolchains()

# #
# # Android
# #
# load("//android:repositories.bzl", "android_repos")
# android_repos()

# load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")
# grpc_java_repositories()

# load("@build_bazel_rules_android//android:sdk_repository.bzl", "android_sdk_repository")
# android_sdk_repository(name = "androidsdk")

# #
# # Closure
# #
# load("//closure:repositories.bzl", "closure_repos")
# closure_repos()

# load("@io_bazel_rules_closure//closure:repositories.bzl", "rules_closure_dependencies", "rules_closure_toolchains")
# rules_closure_dependencies(
#     omit_bazel_skylib = True,
#     omit_com_google_protobuf = True,
#     omit_zlib = True,
# )
# rules_closure_toolchains()

# #
# # Go
# #
# # Load rules_go before running grpc_deps in C++, since that depends on a very old version of
# # rules_go
# #
# load("//:repositories.bzl", "bazel_gazelle", "io_bazel_rules_go")
# io_bazel_rules_go()
# bazel_gazelle()

# load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
# go_rules_dependencies()
# go_register_toolchains()

# load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
# gazelle_dependencies()

# load("//go:repositories.bzl", "go_repos")
# go_repos()

# #
# # C++
# #
# load("//cpp:repositories.bzl", "cpp_repos")
# cpp_repos()

# load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")
# grpc_deps()

# #
# # C#
# #
# load("//csharp:repositories.bzl", "csharp_repos")
# csharp_repos()

# load("@io_bazel_rules_dotnet//dotnet:deps.bzl", "dotnet_repositories")

# dotnet_repositories()

# load(
#     "@io_bazel_rules_dotnet//dotnet:defs.bzl",
#     "core_register_sdk",
#     "net_register_sdk",
#     "dotnet_register_toolchains",
#     "dotnet_repositories_nugets",
# )

# dotnet_register_toolchains()
# dotnet_repositories_nugets()

# core_register_sdk()

# load("@rules_proto_grpc//csharp/nuget:nuget.bzl", "nuget_rules_proto_grpc_packages")
# nuget_rules_proto_grpc_packages()

# #
# # D
# #
# load("//d:repositories.bzl", "d_repos")
# d_repos()

# load("@io_bazel_rules_d//d:d.bzl", "d_repositories")
# d_repositories()

# #
# # Go
# #
# # Moved to above C++

# #
# # gRPC gateway
# #
# load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
# go_rules_dependencies()
# go_register_toolchains()

# load("//github.com/grpc-ecosystem/grpc-gateway:repositories.bzl", "gateway_repos")

# gateway_repos()

# load("@grpc_ecosystem_grpc_gateway//:repositories.bzl", "go_repositories")

# go_repositories()

# #
# # gRPC web
# #
# load("//github.com/grpc/grpc-web:repositories.bzl", "grpc_web_repos")
# grpc_web_repos()

# #
# # Java
# #
# load("//java:repositories.bzl", "java_repos")
# java_repos()

# # grpc_java_repositories already called above in android

# #
# # NodeJS
# #
# load("//nodejs:repositories.bzl", "nodejs_repos")
# nodejs_repos()

# #
# # Objective-C
# #
# load("//objc:repositories.bzl", "objc_repos")
# objc_repos()

# #
# # PHP
# #
# load("//php:repositories.bzl", "php_repos")
# php_repos()

# #
# # Python
# #
# load("//python:repositories.bzl", "python_repos")
# python_repos()

# load("@rules_python//python:repositories.bzl", "py_repositories")
# py_repositories()

# load("@rules_python//python:pip.bzl", "pip_install")
# pip_install(
#     name = "rules_proto_grpc_py3_deps",
#     python_interpreter = "python3",
#     requirements = "//python:requirements.txt",
# )

# #
# # Ruby
# #
# load("//ruby:repositories.bzl", "ruby_repos")
# ruby_repos()

# load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_register_toolchains")
# ruby_register_toolchains()

# load("@com_github_yugui_rules_ruby//ruby/private:bundle.bzl", "bundle_install")
# bundle_install(
#     name = "rules_proto_grpc_gems",
#     gemfile = "//ruby:Gemfile",
#     gemfile_lock = "//ruby:Gemfile.lock",
# )

# #
# # Rust
# #
# load("//rust:repositories.bzl", "rust_repos")
# rust_repos()

# load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")
# rust_repositories()

# load("@io_bazel_rules_rust//:workspace.bzl", "rust_workspace")
# rust_workspace()

# #
# # Scala
# #
# load("//scala:repositories.bzl", "scala_repos")
# scala_repos()

# load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")
# scala_repositories()

# load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")
# scala_register_toolchains()

# load("@io_bazel_rules_scala//scala_proto:scala_proto.bzl", "scala_proto_repositories")
# scala_proto_repositories()

# #
# # Scala routeguide
# #
# load("@bazel_tools//tools/build_defs/repo:jvm.bzl", "jvm_maven_import_external")
# jvm_maven_import_external(
#     name = "com_thesamet_scalapb_scalapb_json4s",
#     artifact = "com.thesamet.scalapb:scalapb-json4s_2.12:0.7.1",
#     artifact_sha256 = "6c8771714329464e03104b6851bfdc3e2e4967276e1a9bd2c87c3b5a6d9c53c7",
#     server_urls = ["https://repo.maven.apache.org/maven2"],
# )

# jvm_maven_import_external(
#     name = "org_json4s_json4s_jackson_2_12",
#     artifact = "org.json4s:json4s-jackson_2.12:3.6.1",
#     artifact_sha256 = "83b854a39e69f022ad3d7dd3da664623252dc822ed4ed1117304f39115c88043",
#     server_urls = ["https://repo.maven.apache.org/maven2"],
# )

# jvm_maven_import_external(
#     name = "org_json4s_json4s_core_2_12",
#     artifact = "org.json4s:json4s-core_2.12:3.6.1",
#     artifact_sha256 = "e0f481509429a24e295b30ba64f567bad95e8d978d0882ec74e6dab291fcdac0",
#     server_urls = ["https://repo.maven.apache.org/maven2"],
# )

# jvm_maven_import_external(
#     name = "org_json4s_json4s_ast_2_12",
#     artifact = "org.json4s:json4s-ast_2.12:3.6.1",
#     artifact_sha256 = "39c7de601df28e32eb0c4e3d684ec65bbf2e59af83c6088cda12688d796f7746",
#     server_urls = ["https://repo.maven.apache.org/maven2"],
# )

# #
# # Swift
# #
# load("//swift:repositories.bzl", "swift_repos")
# swift_repos()

# load(
#     "@build_bazel_rules_swift//swift:repositories.bzl",
#     "swift_rules_dependencies",
# )
# swift_rules_dependencies()

# load(
#     "@build_bazel_apple_support//lib:repositories.bzl",
#     "apple_support_dependencies",
# )
# apple_support_dependencies()

#
# Misc
#
# load("@bazel_gazelle//:deps.bzl", "go_repository")
# go_repository(
#     name = "com_github_urfave_cli",
#     commit = "44cb242eeb4d76cc813fdc69ba5c4b224677e799",
#     importpath = "github.com/urfave/cli",
# )
