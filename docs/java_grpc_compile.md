---
layout: default
title: java_grpc_compile
permalink: java/java_grpc_compile
parent: java
---

# java_grpc_compile

Generates protocol buffer sources for the [java](/java) language.

## `WORKSPACE`

```python
load("@build_stack_rules_proto//toolchains:protoc.bzl", "protoc_toolchain")

protoc_toolchain()

load("@build_stack_rules_proto//rules:java_grpc_compile_deps.bzl", "java_grpc_compile_deps")

java_grpc_compile_deps()

load("@rules_jvm_external//:defs.bzl", "maven_install")
load("@io_grpc_grpc_java//:repositories.bzl", "IO_GRPC_GRPC_JAVA_ARTIFACTS")
load("@io_grpc_grpc_java//:repositories.bzl", "IO_GRPC_GRPC_JAVA_OVERRIDE_TARGETS")

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

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories()
```

## `BUILD.bazel`

```python
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules:java_grpc_compile.bzl", "java_grpc_compile")

proto_library(
    name = "foo_proto",
    srcs = ["foo.proto"],
)

java_grpc_compile(
    name = "java_grpc_compile_foo_proto",
    deps = [":foo_proto"],
)
```

## Plugins

| Label | Tool | Outputs |
| ---- | ---- | ------- |
| `//plugins/java/proto:proto` |  |  |
| `//plugins/java/grpc:grpc` |  |  |

## Dependencies

```python
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def java_grpc_compile_deps():
    bazel_skylib()
    rules_python()
    zlib()
    rules_jvm_external()
    com_google_protobuf()
    io_grpc_grpc_java()
    rules_java()

def bazel_skylib():
    _maybe(
        http_archive,
        name = "bazel_skylib",
        sha256 = "ebdf850bfef28d923a2cc67ddca86355a449b5e4f38b0a70e584dc24e5984aa6",
        strip_prefix = "bazel-skylib-f80bc733d4b9f83d427ce3442be2e07427b2cc8d",
        urls = [
            "https://github.com/bazelbuild/bazel-skylib/archive/f80bc733d4b9f83d427ce3442be2e07427b2cc8d.tar.gz",
        ],
    )

def rules_python():
    _maybe(
        http_archive,
        name = "rules_python",
        sha256 = "8cc0ad31c8fc699a49ad31628273529ef8929ded0a0859a3d841ce711a9a90d5",
        strip_prefix = "rules_python-c7e068d38e2fec1d899e1c150e372f205c220e27",
        urls = [
            "https://github.com/bazelbuild/rules_python/archive/c7e068d38e2fec1d899e1c150e372f205c220e27.tar.gz",
        ],
    )

def zlib():
    _maybe(
        http_archive,
        name = "zlib",
        sha256 = "c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1",
        strip_prefix = "zlib-1.2.11",
        urls = [
            "https://mirror.bazel.build/zlib.net/zlib-1.2.11.tar.gz",
            "https://zlib.net/zlib-1.2.11.tar.gz",
        ],
        build_file = "@build_stack_rules_proto//third_party:BUILD.bazel.zlib",
    )

def rules_jvm_external():
    _maybe(
        http_archive,
        name = "rules_jvm_external",
        sha256 = "cee024d5892c3191937d52909a86cba0ef7b5cdda488d00be84fc37590194339",
        strip_prefix = "rules_jvm_external-576cc9da001be3bae4021ae9e0c06ebb48fcae5d",
        urls = [
            "https://github.com/bazelbuild/rules_jvm_external/archive/576cc9da001be3bae4021ae9e0c06ebb48fcae5d.tar.gz",
        ],
    )

def com_google_protobuf():
    _maybe(
        http_archive,
        name = "com_google_protobuf",
        sha256 = "d0f5f605d0d656007ce6c8b5a82df3037e1d8fe8b121ed42e536f569dec16113",
        strip_prefix = "protobuf-3.14.0",
        urls = [
            "https://github.com/protocolbuffers/protobuf/archive/v3.14.0.tar.gz",
        ],
    )

def io_grpc_grpc_java():
    _maybe(
        http_archive,
        name = "io_grpc_grpc_java",
        sha256 = "82b3cf09f98a5932e1b55175aaec91b2a3f424eec811e47b2a3be533044d9afb",
        strip_prefix = "grpc-java-7f7821c616598ce4e33d2045c5641b2348728cb8",
        urls = [
            "https://github.com/grpc/grpc-java/archive/7f7821c616598ce4e33d2045c5641b2348728cb8.tar.gz",
        ],
    )

def rules_java():
    _maybe(
        http_archive,
        name = "rules_java",
        sha256 = "7c4bbe11e41c61212a5cf16d9aafaddade3f5b1b6c8bf94270d78215fafd4007",
        strip_prefix = "rules_java-c13e3ead84afb95f81fbddfade2749d8ba7cb77f",
        urls = [
            "https://github.com/bazelbuild/rules_java/archive/c13e3ead84afb95f81fbddfade2749d8ba7cb77f.tar.gz",
        ],
    )

```