---
layout: default
title: Installation
permalink: guides/installation
parent: Guides
nav_order: 1
---

Getting Started with rules_proto involves 2 steps:

1. Setup the `WORKSPACE` file.
1. Setup `BUILD.bazel` in the root of the workspace directory.

## `./WORKSPACE`

```python
# == rules_proto ==

http_archive(
    name = "build_stack_rules_proto",
    sha256 = "c70122c7c5213a7492ae6e5274925505e13564634bb3541b92ac76c4059760a8",
    strip_prefix = "rules_proto-fe9ef7534f97ec7a42d987ab49835505b1cec56e",
    urls = ["https://github.com/stackb/rules_proto/archive/fe9ef7534f97ec7a42d987ab49835505b1cec56e.tar.gz"],
)

register_toolchains("@build_stack_rules_proto//toolchain")

load("@build_stack_rules_proto//deps:core_deps.bzl", "core_deps")

core_deps()

# == Go ==

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.18.2")

# == Gazelle ==

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()

# == Protobuf ==

load("@build_stack_rules_proto//deps:protobuf_core_deps.bzl", "protobuf_core_deps")

protobuf_core_deps()
```

> NOTE: If you already have `rules_go` and `bazel_gazelle` in your workspace,
> you can skip the `core_deps`.  If you already have `com_google_protobuf`, you
> can skip `protobuf_core_deps`.

## `./BUILD.bazel`

The simplest configuration would be:

```python
load("@bazel_gazelle//:def.bzl", "gazelle")

gazelle(
    name = "update_build_files",
    gazelle = "@build_stack_rules_proto//:gazelle",
)
```

If you have your own gazelle language implementation(s), use something like the following:

```python
load("@bazel_gazelle//:def.bzl", "gazelle", "gazelle_binary")

gazelle_binary(
    name = "gazelle",
    languages = [
        # NOTE: order matters here
        "@bazel_gazelle//language/proto",
        "@build_stack_rules_proto//language/protobuf",
        "@bazel_gazelle//language/go",
        "//gazelle/language/foo",
    ],
)

gazelle(
    name = "update_build_files",
    gazelle = ":gazelle",
)
```

You should now be able to `bazel run //:update_build_files`.  Proceed to [gazelle configuration](configuration).