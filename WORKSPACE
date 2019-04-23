workspace(name = "build_stack_rules_proto")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

# local_repository(
#     name = "org_pubref_rules_node",
#     path = "/home/pcj/github/pubref/rules_node",
# )

# **************************************************************
#
#
# cpp
#
# **************************************************************

load("@build_stack_rules_proto//cpp:deps.bzl", "cpp_grpc_library")

cpp_grpc_library()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()


# **************************************************************
#
#
# closure
#
# **************************************************************

load("@build_stack_rules_proto//closure:deps.bzl", "closure_proto_library")

closure_proto_library()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories()


# **************************************************************
#
#
# csharp
#
# **************************************************************

load("@build_stack_rules_proto//csharp:deps.bzl", "csharp_grpc_library")

csharp_grpc_library()

load("@io_bazel_rules_dotnet//dotnet:defs.bzl", 
    "dotnet_register_toolchains", 
    "core_register_sdk",
    "dotnet_repositories",
)

core_version = "v2.1.503"

dotnet_register_toolchains(
    core_version = core_version,
)

dotnet_register_toolchains(
    core_version = core_version,
)

core_register_sdk(
    name = "core_sdk",
    core_version = core_version
)

dotnet_repositories()

load("@build_stack_rules_proto//csharp/nuget:packages.bzl", nuget_packages = "packages")

nuget_packages()

load("@build_stack_rules_proto//csharp/nuget:nuget.bzl", "nuget_protobuf_packages")

nuget_protobuf_packages()

load("@build_stack_rules_proto//csharp/nuget:nuget.bzl", "nuget_grpc_packages")

nuget_grpc_packages()


# **************************************************************
#
#
# go
#
# **************************************************************

load("@build_stack_rules_proto//go:deps.bzl", "go_grpc_library")

go_grpc_library()

load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()


# **************************************************************
#
#
# java
#
# **************************************************************

load("@build_stack_rules_proto//:deps.bzl", "io_grpc_grpc_java")

io_grpc_grpc_java()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories(omit_com_google_protobuf = True)

load("@build_stack_rules_proto//java:deps.bzl", "java_grpc_library")

java_grpc_library()


# **************************************************************
#
#
# node
#
# **************************************************************

load("@build_stack_rules_proto//node:deps.bzl", "node_grpc_library")

node_grpc_library()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@org_pubref_rules_node//node:rules.bzl", "node_repositories", "yarn_modules")

node_repositories()

yarn_modules(
    name = "proto_node_modules",
    deps = {
        "google-protobuf": "3.6.1",
    },
)

yarn_modules(
    name = "grpc_node_modules",
    deps = {
        "grpc": "1.15.1",
        "async": "2.6.1",
    },
)


# **************************************************************
#
#
# python
#
# **************************************************************

load("@build_stack_rules_proto//python:deps.bzl", "python_grpc_library")

python_grpc_library()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@io_bazel_rules_python//python:pip.bzl", "pip_import", "pip_repositories")

pip_repositories()

pip_import(
    name = "protobuf_py_deps",
    requirements = "@build_stack_rules_proto//python/requirements:protobuf.txt",
)

load("@protobuf_py_deps//:requirements.bzl", protobuf_pip_install = "pip_install")

protobuf_pip_install()

pip_import(
    name = "grpc_py_deps",
    requirements = "@build_stack_rules_proto//python:requirements.txt",
)

load("@grpc_py_deps//:requirements.bzl", grpc_pip_install = "pip_install")

grpc_pip_install()


# **************************************************************
#
#
# scala
#
# **************************************************************

load("@build_stack_rules_proto//scala:deps.bzl", "scala_grpc_library")

scala_grpc_library()

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")

scala_repositories()

load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")

scala_register_toolchains()

load("@io_bazel_rules_scala//scala_proto:scala_proto.bzl", "scala_proto_repositories")

scala_proto_repositories()


# **************************************************************
#
#
# swift
#
# **************************************************************

load("@build_stack_rules_proto//swift:deps.bzl", "swift_grpc_library")

swift_grpc_library()

load(
    "@build_bazel_rules_swift//swift:repositories.bzl",
    "swift_rules_dependencies",
)

swift_rules_dependencies()

load(
    "@build_bazel_apple_support//lib:repositories.bzl",
    "apple_support_dependencies",
)

apple_support_dependencies()


# **************************************************************
#
#
# ruby
#
# **************************************************************

load("//ruby:deps.bzl", "ruby_grpc_library")

ruby_grpc_library()

load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_register_toolchains")

ruby_register_toolchains()

load("@com_github_yugui_rules_ruby//ruby/private:bundle.bzl", "bundle_install")

bundle_install(
    name = "routeguide_gems_bundle",
    gemfile = "//ruby:Gemfile",
    gemfile_lock = "//ruby:Gemfile.lock",
)


# **************************************************************
#
#
# dart
#
# **************************************************************

load("//dart:deps.bzl", "dart_grpc_library")

dart_grpc_library()

load("@io_bazel_rules_dart//dart/build_rules:repositories.bzl", "dart_repositories")

dart_repositories()

load("@dart_pub_deps_protoc_plugin//:deps.bzl", dart_protoc_plugin_deps = "pub_deps")

dart_protoc_plugin_deps()

load("@dart_pub_deps_grpc//:deps.bzl", dart_grpc_deps = "pub_deps")

dart_grpc_deps()



# **************************************************************
#
#
# gazelle & buildifier
#
# **************************************************************

load("//:deps.bzl", "bazel_gazelle", "com_github_bazelbuild_buildtools")

com_github_bazelbuild_buildtools()

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")

buildifier_dependencies()

bazel_gazelle()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

# gazelle:repo bazel_gazelle


# **************************************************************
#
#
# rust
#
# **************************************************************

load("//rust:deps.bzl", "rust_grpc_library")

rust_grpc_library()

load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")

rust_repositories()

load("@io_bazel_rules_rust//:workspace.bzl", "bazel_version")

bazel_version(name = "bazel_version")

load("@io_bazel_rules_rust//proto/raze:crates.bzl", "raze_fetch_remote_crates")

raze_fetch_remote_crates()

# =========================================

# PROTOBUF_VERSION = "3.6.1.3"

# http_archive(
#     name = "com_google_protobuf",
#     sha256 = "73fdad358857e120fd0fa19e071a96e15c0f23bb25f85d3f7009abfd4f264a2a",
#     strip_prefix = "protobuf-" + PROTOBUF_VERSION,
#     url = "https://github.com/google/protobuf/archive/v%s.tar.gz" % PROTOBUF_VERSION,
# )

# =========================================
# load("//:deps.bzl", "com_github_grpc_grpc")

# com_github_grpc_grpc()

# load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", com_github_grpc_grpc_bazel_grpc_deps = "grpc_deps")

# com_github_grpc_grpc_bazel_grpc_deps()

# # =========================================
# load("//:deps.bzl", "io_grpc_grpc_java")

# io_grpc_grpc_java()

# load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

# grpc_java_repositories(
#     omit_com_google_protobuf = True,
# )

# # =========================================

# load("//csharp:deps.bzl", "csharp_grpc_library")

# csharp_grpc_library()

# load("@io_bazel_rules_dotnet//dotnet:defs.bzl", 
#     "dotnet_register_toolchains", 
#     "core_register_sdk",
#     "dotnet_repositories",
# )

# core_version = "v2.1.503"

# dotnet_register_toolchains(
#     core_version = core_version,
# )

# dotnet_register_toolchains(
#     core_version = core_version,
# )

# core_register_sdk(
#     name = "core_sdk",
#     core_version = core_version
# )

# dotnet_repositories()

# load("//csharp/nuget:packages.bzl", nuget_packages = "packages")

# nuget_packages()

# load("//csharp/nuget:nuget.bzl", "nuget_protobuf_packages")

# nuget_protobuf_packages()

# load("//csharp/nuget:nuget.bzl", "nuget_grpc_packages")

# nuget_grpc_packages()

# # =========================================

# load("//python:deps.bzl", "python_grpc_library")

# python_grpc_library()

# load("@io_bazel_rules_python//python:pip.bzl", "pip_import", "pip_repositories")

# pip_repositories()

# pip_import(
#     name = "protobuf_py_deps",
#     requirements = "@build_stack_rules_proto//python/requirements:protobuf.txt",
# )

# load("@protobuf_py_deps//:requirements.bzl", protobuf_pip_install = "pip_install")

# protobuf_pip_install()

# pip_import(
#     name = "grpc_py_deps",
#     requirements = "@build_stack_rules_proto//python:requirements.txt",
# )

# load("@grpc_py_deps//:requirements.bzl", grpc_pip_install = "pip_install")

# grpc_pip_install()

# # =========================================

# load("//java:deps.bzl", "java_grpc_library")

# java_grpc_library()

# # =========================================

# load("//scala:deps.bzl", "scala_grpc_library")

# scala_grpc_library()

# load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")

# scala_repositories()

# load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")

# scala_register_toolchains()

# load("@io_bazel_rules_scala//scala_proto:scala_proto.bzl", "scala_proto_repositories")

# scala_proto_repositories()

# # See rules_scala WORKSPACE
# register_toolchains("@io_bazel_rules_scala//test/proto:scalapb_toolchain")

# # ===========

# load("//closure:deps.bzl", "closure_proto_library")

# closure_proto_library()

# load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

# closure_repositories(
#     omit_com_google_protobuf = True,
# )

# # =========================================

# load("//node:deps.bzl", "node_grpc_library")

# node_grpc_library()

# load("@org_pubref_rules_node//node:rules.bzl", "node_repositories", "yarn_modules")

# node_repositories()

# yarn_modules(
#     name = "proto_node_modules",
#     deps = {
#         "google-protobuf": "3.6.1",
#     },
# )

# yarn_modules(
#     name = "grpc_node_modules",
#     deps = {
#         "grpc": "1.15.1",
#         "async": "2.6.1",
#     },
# )

# # =========================================

# load("//go:deps.bzl", "go_grpc_library")

# go_grpc_library()

# load("@io_bazel_rules_go//go:def.bzl", "go_register_toolchains", "go_rules_dependencies")

# go_rules_dependencies()

# go_register_toolchains()

# load("//:deps.bzl", "bazel_gazelle")

# bazel_gazelle()

# load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

# gazelle_dependencies()

# # gazelle:repo bazel_gazelle

# # =========================================

# load("//dart:deps.bzl", "dart_grpc_library")

# dart_grpc_library()

# load("@io_bazel_rules_dart//dart/build_rules:repositories.bzl", "dart_repositories")

# dart_repositories()

# load("@dart_pub_deps_protoc_plugin//:deps.bzl", dart_protoc_plugin_deps = "pub_deps")

# dart_protoc_plugin_deps()

# load("@dart_pub_deps_grpc//:deps.bzl", dart_grpc_deps = "pub_deps")

# dart_grpc_deps()

# # =========================================

# load("//github.com/stackb/grpc.js:deps.bzl", "closure_grpc_library")

# closure_grpc_library()

# # =========================================

# load("//github.com/grpc/grpc-web:deps.bzl", "closure_grpc_library")

# closure_grpc_library()

# # =========================================

# load("//github.com/gogo/protobuf:deps.bzl", "gogo_grpc_library")

# gogo_grpc_library()

# # =========================================

# ##
# ## Must have swiftc on your system.  To install on ubuntu 17.10, here's what I did:
# ##
# ## $ sudo apt install ubuntu-make; umake swift
# ## $ vi `/.bashrc` and update PATH to include ~/.local/share/umake/swift/swift-lang/usr/bin
# ##

# # load("//swift:repositories.bzl", "swift_toolchain")

# # # Default values work with linux, x86_64, /usr/local/bin/clang.
# # swift_toolchain(
# #     root = "/home/pcj/.local/share/umake/swift/swift-lang/usr",
# # )

# load("//swift:deps.bzl", "swift_grpc_library")

# swift_grpc_library()

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


# # =======================================

# # load("//ruby:deps.bzl", "ruby_grpc_library")

# # ruby_grpc_library()

# # load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_register_toolchains")

# # ruby_register_toolchains()

# # load("@com_github_yugui_rules_ruby//ruby/private:bundle.bzl", "bundle_install")

# # bundle_install(
# #     name = "routeguide_gems_bundle",
# #     gemfile = "//ruby:Gemfile",
# #     gemfile_lock = "//ruby:Gemfile.lock",
# # )

# # =======================================

# # load("//rust:deps.bzl", "rust_grpc_library")

# # rust_grpc_library()

# # load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")

# # rust_repositories()

# # load("//rust/cargo:crates.bzl", "raze_fetch_remote_crates")

# # raze_fetch_remote_crates()

# # ========================================

# load("//android:deps.bzl", "android_grpc_library")

# android_grpc_library()

# load("@build_bazel_rules_android//android:sdk_repository.bzl", "android_sdk_repository")

# #
# # This version of android_sdk_repository uses the ANDROID_HOME env var.
# #
# # Why is this not a workspace rule?  This is so painful.
# # $ curl -O -J -L https://dl.google.com/android/repository/sdk-tools-linux-4333796.zip
# # $ mkdir ~/android-sdk
# # $ unzip in that dir
# # $ tools/bin/sdkmanager --list
# # $ tools/bin/sdkmanager "platforms;android-28" "build-tools;28.0.3"
# # $ export ANDROID_HOME="/home/pcj/android-sdk/"
# #
# android_sdk_repository(
#     name = "androidsdk",
# )

# load("@gmaven_rules//:gmaven.bzl", "gmaven_rules")

# gmaven_rules()

# # =======================================

# load("//github.com/grpc-ecosystem/grpc-gateway:deps.bzl", "gateway_grpc_library")

# gateway_grpc_library()

# # load("@grpc_ecosystem_grpc_gateway//:repositories.bzl", grpc_gateway_repositories = "repositories")

# # grpc_gateway_repositories()

# load("//tools/exgen:deps.bzl", "exgen_deps")

# exgen_deps()

# # =======================================

# http_archive(
#     name = "com_github_bazelbuild_buildtools",
#     strip_prefix = "buildtools-41d89cd7c8328bb912f3b8f50d2dc970805d21f8",
#     url = "https://github.com/bazelbuild/buildtools/archive/41d89cd7c8328bb912f3b8f50d2dc970805d21f8.zip",
# )

# load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")

# buildifier_dependencies()




# **************************************************************
#
#
# tools & other misc support (not language specific)
#
# **************************************************************

go_repository(
    name = "com_github_urfave_cli",
    commit = "44cb242eeb4d76cc813fdc69ba5c4b224677e799",
    importpath = "github.com/urfave/cli",
)
