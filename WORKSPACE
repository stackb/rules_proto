workspace(name = "build_stack_rules_proto")

#
# Core
#
load("//protobuf:deps.bzl", "protobuf_deps")
protobuf_deps()


#
# Android
#
load("//android:deps.bzl", "android_deps")
android_deps()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")
grpc_java_repositories(
    omit_com_google_protobuf = True,
    omit_net_zlib = True
)

load("@build_bazel_rules_android//android:sdk_repository.bzl", "android_sdk_repository")
android_sdk_repository(name = "androidsdk")

#
# Android routeguide
#
load("//:deps.bzl", "rules_jvm_external")
rules_jvm_external()

load("@rules_jvm_external//:defs.bzl", "maven_install")
maven_install(
    name = "maven_android",
    artifacts = [
        "com.android.support:appcompat-v7:28.0.0",
    ],
    # Fail if a checksum file for the artifact is missing in the repository.
    # Falls through "SHA-1" and "MD5". Defaults to True.
    fail_on_missing_checksum = False,
    repositories = [
        "https://maven.google.com",
        "https://repo1.maven.org/maven2",
    ],
)


#
# Closure
#
load("//closure:deps.bzl", "closure_deps")
closure_deps()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")
closure_repositories()


#
# C++
#
load("//cpp:deps.bzl", "cpp_deps")
cpp_deps()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")
grpc_deps()


#
# csharp
#
load("//csharp:deps.bzl", "csharp_deps")
csharp_deps()

load(
    "@io_bazel_rules_dotnet//dotnet:defs.bzl",
    "core_register_sdk",
    "dotnet_register_toolchains",
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
    name = "core_sdk_{}".format(core_version),
    core_version = core_version,
)

dotnet_repositories()

load("//csharp/nuget:packages.bzl", nuget_packages = "packages")
nuget_packages()

load("//csharp/nuget:nuget.bzl", "nuget_protobuf_packages")
nuget_protobuf_packages()

load("//csharp/nuget:nuget.bzl", "nuget_grpc_packages")
nuget_grpc_packages()


#
# D
#
load("//d:deps.bzl", "d_deps")
d_deps()

load("@io_bazel_rules_d//d:d.bzl", "d_repositories")
d_repositories()


#
# Go
#
load("//go:deps.bzl", "go_deps")
go_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
go_rules_dependencies()
go_register_toolchains()


#
# grpc.js
#
load("//github.com/stackb/grpc.js:deps.bzl", "grpcjs_deps")
grpcjs_deps()


#
# gRPC gateway
#
load("//:deps.bzl", "bazel_gazelle", "io_bazel_rules_go")
io_bazel_rules_go()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
go_rules_dependencies()
go_register_toolchains()
bazel_gazelle()

load("//github.com/grpc-ecosystem/grpc-gateway:deps.bzl", "gateway_deps")
gateway_deps()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()


#
# gRPC web
#
load("//github.com/grpc/grpc-web:deps.bzl", "grpc_web_deps")
grpc_web_deps()


#
# Java
#
load("//java:deps.bzl", "java_deps")
java_deps()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")
grpc_java_repositories(omit_com_google_protobuf = True, omit_net_zlib = True)


#
# NodeJS
#
load("//nodejs:deps.bzl", "nodejs_deps")
nodejs_deps()

load("@build_bazel_rules_nodejs//:defs.bzl", "yarn_install")
yarn_install(
    name = "nodejs_modules",
    package_json = "//nodejs:requirements/package.json",
    yarn_lock = "//nodejs:requirements/yarn.lock",
)


#
# Objective-C
#
load("//objc:deps.bzl", "objc_deps")
objc_deps()


#
# PHP
#
load("//php:deps.bzl", "php_deps")
php_deps()


#
# Python
#
load("//python:deps.bzl", "python_deps")
python_deps()

load("@com_apt_itude_rules_pip//rules:dependencies.bzl", "pip_rules_dependencies")
pip_rules_dependencies()

load("@com_apt_itude_rules_pip//rules:repository.bzl", "pip_repository")
pip_repository(
    name = "grpc_py2_deps",
    python_interpreter = "python2",
    requirements = "//python/requirements:grpc.txt",
)
pip_repository(
    name = "grpc_py3_deps",
    python_interpreter = "python3",
    requirements = "//python/requirements:grpc.txt",
)


#
# Ruby
#
load("//ruby:deps.bzl", "ruby_deps")
ruby_deps()

load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_register_toolchains")
ruby_register_toolchains()

load("@com_github_yugui_rules_ruby//ruby/private:bundle.bzl", "bundle_install")
bundle_install(
    name = "routeguide_gems_bundle",
    gemfile = "//ruby:Gemfile",
    gemfile_lock = "//ruby:Gemfile.lock",
)


#
# Rust
#
load("//rust:deps.bzl", "rust_deps")
rust_deps()

load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")
rust_repositories()

load("@io_bazel_rules_rust//:workspace.bzl", "bazel_version")
bazel_version(name = "bazel_version")

load("@io_bazel_rules_rust//proto:repositories.bzl", "rust_proto_repositories")
rust_proto_repositories()


#
# Scala
#
load("//scala:deps.bzl", "scala_deps")
scala_deps()

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")
scala_repositories()

load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")
scala_register_toolchains()

load("@io_bazel_rules_scala//scala_proto:scala_proto.bzl", "scala_proto_repositories")
scala_proto_repositories()

#
# Scala routeguide
#
load("@bazel_tools//tools/build_defs/repo:jvm.bzl", "jvm_maven_import_external")
jvm_maven_import_external(
    name = "com_thesamet_scalapb_scalapb_json4s",
    artifact = "com.thesamet.scalapb:scalapb-json4s_2.12:0.7.1",
    server_urls = ["http://central.maven.org/maven2"],
    artifact_sha256 = "6c8771714329464e03104b6851bfdc3e2e4967276e1a9bd2c87c3b5a6d9c53c7",
)

jvm_maven_import_external(
    name = "org_json4s_json4s_jackson_2_12",
    artifact = "org.json4s:json4s-jackson_2.12:3.6.1",
    server_urls = ["http://central.maven.org/maven2"],
    artifact_sha256 = "83b854a39e69f022ad3d7dd3da664623252dc822ed4ed1117304f39115c88043",
)

jvm_maven_import_external(
    name = "org_json4s_json4s_core_2_12",
    artifact = "org.json4s:json4s-core_2.12:3.6.1",
    server_urls = ["http://central.maven.org/maven2"],
    artifact_sha256 = "e0f481509429a24e295b30ba64f567bad95e8d978d0882ec74e6dab291fcdac0",
)

jvm_maven_import_external(
    name = "org_json4s_json4s_ast_2_12",
    artifact = "org.json4s:json4s-ast_2.12:3.6.1",
    server_urls = ["http://central.maven.org/maven2"],
    artifact_sha256 = "39c7de601df28e32eb0c4e3d684ec65bbf2e59af83c6088cda12688d796f7746",
)


#
# Swift
#
load("//swift:deps.bzl", "swift_deps")
swift_deps()

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


#
# Misc
#
load("//:deps.bzl", "bazel_gazelle", "com_github_bazelbuild_buildtools")
com_github_bazelbuild_buildtools()
bazel_gazelle()

load("@com_github_bazelbuild_buildtools//buildifier:deps.bzl", "buildifier_dependencies")
buildifier_dependencies()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")
gazelle_dependencies()

go_repository(
    name = "com_github_urfave_cli",
    commit = "44cb242eeb4d76cc813fdc69ba5c4b224677e799",
    importpath = "github.com/urfave/cli",
)
