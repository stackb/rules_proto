---
layout: default
title: nodejs_proto_library
permalink: nodejs/nodejs_proto_library
parent: nodejs
---

# nodejs_proto_library

Generates protocol buffer sources for the [nodejs](/rules_proto/nodejs) language.

## `WORKSPACE`

```python
load("@build_stack_rules_proto//toolchains:protoc.bzl", "protoc_toolchain")

protoc_toolchain()

load("@build_stack_rules_proto//rules:nodejs_proto_library_deps.bzl", "nodejs_proto_library_deps")

nodejs_proto_library_deps()

load("@build_bazel_rules_nodejs//:index.bzl", "yarn_install")

yarn_install(
    name = "google_protobuf_node_modules",
    package_json = "@build_stack_rules_proto//plugins/nodejs/modules/google-protobuf:package.json",
    symlink_node_modules = False,
    yarn_lock = "@build_stack_rules_proto//plugins/nodejs/modules/google-protobuf:yarn.lock",
)
```

## `BUILD.bazel`

```python
load("@rules_proto//proto:defs.bzl", "proto_library")
load("@build_stack_rules_proto//rules:nodejs_proto_library.bzl", "nodejs_proto_library")

proto_library(
    name = "foo_proto",
    srcs = ["foo.proto"],
)

nodejs_proto_library(
    name = "nodejs_proto_library_foo_proto",
    deps = [":foo_proto"],
)
```

## Plugins

| Label | Tool | Outputs |
| ---- | ---- | ------- |
| `//plugins/nodejs/proto:proto` |  |  `{protopath}_pb.js` |

## Dependencies

```python
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive", "http_file")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def nodejs_proto_library_deps():
    build_bazel_rules_nodejs()  # via google_protobuf_node_modules
    zlib()  # via com_google_protobuf
    rules_python()  # via com_google_protobuf
    bazel_skylib()  # via com_google_protobuf
    com_google_protobuf()  # via rule nodejs_proto_library



def build_bazel_rules_nodejs():
    _maybe(
        http_archive,
        name = "build_bazel_rules_nodejs",
        sha256 = "6142e9586162b179fdd570a55e50d1332e7d9c030efd853453438d607569721d",
        urls = [
            "https://github.com/bazelbuild/rules_nodejs/releases/download/3.0.0/rules_nodejs-3.0.0.tar.gz",
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
```