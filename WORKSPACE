workspace(name = "build_stack_rules_proto")

# gazelle:repo bazel_gazelle

# ----------------------------------------------------
# Toolchain-Related
# ----------------------------------------------------

register_toolchains("//toolchain:standard")

# ----------------------------------------------------
# Top-Level Dependency Trees
# ----------------------------------------------------

load("//deps:core_deps.bzl", "core_deps")

core_deps()

load("//deps:protobuf_core_deps.bzl", "protobuf_core_deps")

protobuf_core_deps()

load("//deps:prebuilt_protoc_deps.bzl", "prebuilt_protoc_deps")

prebuilt_protoc_deps()

load("//deps:grpc_core_deps.bzl", "grpc_core_deps")

grpc_core_deps()

load("//deps:grpc_java_deps.bzl", "grpc_java_deps")

grpc_java_deps()

load("//deps:closure_deps.bzl", "closure_deps")

closure_deps()

load("//deps:grpc_js_deps.bzl", "grpc_js_deps")

grpc_js_deps()

load("//deps:scala_deps.bzl", "scala_deps")

scala_deps()

load("//deps:nodejs_deps.bzl", "nodejs_deps")

nodejs_deps()

load("//deps:grpc_node_deps.bzl", "grpc_node_deps")

grpc_node_deps()

load("//deps:grpc_web_deps.bzl", "grpc_web_deps")

grpc_web_deps()

load("//deps:ts_proto_deps.bzl", "ts_proto_deps")

ts_proto_deps()

load("//deps:example_routeguide_nodejs_deps.bzl", "example_routeguide_nodejs_deps")

example_routeguide_nodejs_deps()

# ----------------------------------------------------
# Go Tools
# ----------------------------------------------------

load(
    "@io_bazel_rules_go//go:deps.bzl",
    "go_register_toolchains",
    "go_rules_dependencies",
)

go_rules_dependencies()

go_register_toolchains(version = "1.18.2")

# ----------------------------------------------------
# Gazelle
# ----------------------------------------------------

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

load("//:go_deps.bzl", "go_deps")

# gazelle:repository_macro go_deps.bzl%go_deps
go_deps()

# ----------------------------------------------------
# Core gRPC
# ----------------------------------------------------

load(
    "@com_github_grpc_grpc//bazel:grpc_deps.bzl",
    "grpc_deps",
)

grpc_deps()

# ----------------------------------------------------
# Java
# ----------------------------------------------------

load(
    "@rules_jvm_external//:defs.bzl",
    "maven_install",
)
load(
    "@io_grpc_grpc_java//:repositories.bzl",
    "IO_GRPC_GRPC_JAVA_ARTIFACTS",
    "IO_GRPC_GRPC_JAVA_OVERRIDE_TARGETS",
    "grpc_java_repositories",
)

maven_install(
    artifacts = IO_GRPC_GRPC_JAVA_ARTIFACTS,
    generate_compat_repositories = True,
    override_targets = IO_GRPC_GRPC_JAVA_OVERRIDE_TARGETS,
    repositories = [
        "https://repo.maven.apache.org/maven2/",
    ],
)

load(
    "@maven//:compat.bzl",
    "compat_repositories",
)

compat_repositories()

grpc_java_repositories()

# ----------------------------------------------------
# Golang
# ----------------------------------------------------

load("//deps:go_core_deps.bzl", "go_core_deps")

go_core_deps()

# ----------------------------------------------------
# Scala
# ----------------------------------------------------

load("@io_bazel_rules_scala//:scala_config.bzl", "scala_config")

scala_config(scala_version = "2.12.11")

load("@io_bazel_rules_scala//scala:scala.bzl", "scala_repositories")

scala_repositories()

load("@io_bazel_rules_scala//scala:toolchains.bzl", "scala_register_toolchains")

scala_register_toolchains()

# bazel run @maven_scala//:pin, but first comment out the "maven_install_json"
# (put it back once pinned again)
maven_install(
    name = "maven_scala",
    artifacts = [
        "com.thesamet.scalapb:lenses_2.12:0.11.5",
        "com.thesamet.scalapb:scalapb-json4s_2.12:0.12.0",
        "com.thesamet.scalapb:scalapb-runtime_2.12:0.11.5",
        "com.thesamet.scalapb:scalapb-runtime-grpc_2.12:0.11.5",
        "com.thesamet.scalapb:scalapbc_2.12:0.11.5",
        "io.grpc:grpc-api:1.40.1",
        "io.grpc:grpc-core:1.40.1",
        "io.grpc:grpc-netty:1.40.1",
        "io.grpc:grpc-protobuf:1.40.1",
        "io.grpc:grpc-stub:1.40.1",
        "org.json4s:json4s-core_2.12:4.0.3",
    ],
    fetch_sources = True,
    maven_install_json = "//:maven_scala_install.json",
    repositories = ["https://repo1.maven.org/maven2"],
)

load("@maven_scala//:defs.bzl", pinned_maven_scala_install = "pinned_maven_install")

pinned_maven_scala_install()

# ----------------------------------------------------
# Scala/Akka
# ----------------------------------------------------

# bazel run @maven_akka//:pin, but first comment out the "maven_install_json"
# (put it back once pinned again)
maven_install(
    name = "maven_akka",
    artifacts = [
        "com.lightbend.akka.grpc:akka-grpc-codegen_2.12:2.1.3",
        "com.lightbend.akka.grpc:akka-grpc-runtime_2.12:2.1.3",
    ],
    fetch_sources = True,
    maven_install_json = "//:maven_akka_install.json",
    repositories = ["https://repo1.maven.org/maven2"],
)

load("@maven_akka//:defs.bzl", pinned_maven_akka_install = "pinned_maven_install")

pinned_maven_akka_install()

# ----------------------------------------------------
# Closure
# ----------------------------------------------------

load("@io_bazel_rules_closure//closure:repositories.bzl", "rules_closure_dependencies", "rules_closure_toolchains")

rules_closure_toolchains()

rules_closure_dependencies()

# ----------------------------------------------------
# NodeJS
# ----------------------------------------------------

load("@build_bazel_rules_nodejs//:index.bzl", "node_repositories")

node_repositories()

register_toolchains("//toolchain:nodejs")
