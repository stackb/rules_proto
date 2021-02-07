---
layout: default
title: cc_proto_compile
---

# cc_proto_compile

## `WORKSPACE`

```python
load("@build_stack_rules_proto//cc:cc_proto_compile_deps.bzl", "cc_proto_compile_deps")

cc_proto_compile_deps()
```

## `BUILD.bazel`

```python
load("@build_stack_rules_proto//cc:cc_proto_compile.bzl", "cc_proto_compile")

cc_proto_compile(
    name = "cc_proto_compile_person_proto",
    deps = ["@build_stack_rules_proto//example/proto/v1:person_proto"],
)
```

## Plugins

| Label | Tool | Outputs |
| ---- | ---- | ------- |
| `//cc:cc_plugin` |  |  `{protopath}.pb.h` `{protopath}.pb.cc` |


## Dependencies

```python
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def cc_proto_compile_deps():


```