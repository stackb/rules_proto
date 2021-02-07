---
layout: default
title: py_proto_compile
permalink: python/py_proto_compile
parent: python
---

# py_proto_compile

Generates protocol buffer source code for the **python**.

## `WORKSPACE`

```python
load("@build_stack_rules_proto//python:py_proto_compile_deps.bzl", "py_proto_compile_deps")

py_proto_compile_deps()
```

## `BUILD.bazel`

```python
load("@build_stack_rules_proto//python:py_proto_compile.bzl", "py_proto_compile")

py_proto_compile(
    name = "py_proto_compile_person_proto",
    deps = ["@build_stack_rules_proto//example/proto/v1:person_proto"],
)
```

## Plugins

| Label | Tool | Outputs |
| ---- | ---- | ------- |
| `//python:python_plugin` |  |  `{protopath|python}_pb2.py` |


## Dependencies

```python
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def py_proto_compile_deps():
    com_google_protobuf()

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