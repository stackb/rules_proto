---
layout: default
title: Home
permalink: /
nav_order: 0
---

# rules_proto

<table border="0" style="text-align: center"><tr>
<td><img src="https://bazel.build/images/bazel-icon.svg" style="height: 160px"/></td>
<td><img src="data:image/svg+xml,%3C!--%20Created%20by%20Jose%20Rivera%20--%3E%0A%3Csvg%20xmlns%3D%22http%3A%2F%2Fwww.w3.org%2F2000%2Fsvg%22%20xmlns%3Axlink%3D%22http%3A%2F%2Fwww.w3.org%2F1999%2Fxlink%22%20version%3D%221.1%22%20x%3D%220px%22%20y%3D%220px%22%20viewBox%3D%220%200%20100%20100%22%20style%3D%22enable-background%3Anew%200%200%20100%20100%3B%22%20xml%3Aspace%3D%22preserve%22%3E%0A%09%3Cg%3E%0A%09%09%3Cg%3E%0A%09%09%09%3Cpath%20d%3D%22M68.9%2C52.4L68.9%2C52.4l5.3%2C3.4L68.9%2C52.4z%22%20%2F%3E%0A%09%09%3C%2Fg%3E%0A%09%09%3Cpath%20d%3D%22M90.9%2C69.5l0-0.2l0%2C0.8L90.9%2C69.5z%22%20%2F%3E%0A%09%09%3Cpath%20d%3D%22M16.2%2C34.4L15.7%2C34l-0.2-0.2v0l0.2%2C0.1L16.2%2C34.4L16.2%2C34.4L16.2%2C34.4z%20M16.2%2C34.4L15.7%2C34l-0.2-0.2v0l0.2%2C0.1L16.2%2C34.4%20%20%20L16.2%2C34.4L16.2%2C34.4z%20M15.2%2C33.5L15.2%2C33.5l0.2%2C0.2L15.2%2C33.5z%20M16.2%2C34.4L15.7%2C34l-0.2-0.2v0l0.2%2C0.1L16.2%2C34.4L16.2%2C34.4%20%20%20L16.2%2C34.4z%20M16.2%2C34.4L15.7%2C34h0L16.2%2C34.4L16.2%2C34.4L16.2%2C34.4z%22%20%2F%3E%0A%09%09%3Cg%3E%0A%09%09%09%3Cpath%20d%3D%22M16.2%2C34.4L15.7%2C34l-0.2-0.2v0l0.2%2C0.1L16.2%2C34.4L16.2%2C34.4L16.2%2C34.4z%22%20%2F%3E%0A%09%09%3C%2Fg%3E%0A%09%09%3Cpath%20d%3D%22M52.4%2C31.2L52.4%2C31.2L52.4%2C31.2L52.4%2C31.2z%22%20%2F%3E%0A%09%09%3Cg%3E%0A%09%09%09%3Cpath%20d%3D%22M68%2C32.8L68%2C32.8l0.7%2C0.6L68%2C32.8z%22%20%2F%3E%0A%09%09%3C%2Fg%3E%0A%09%09%3Cg%3E%0A%09%09%09%3Cpath%20d%3D%22M94.7%2C82.9l-2.9%2C2.9l-0.9-15.8l0-0.8l-0.2-0.2l-10.1-9.1l4.6-6.4l-0.9-6.7l-9.5-8.3l-0.6-0.5l-5.5-4.7L68%2C32.8v0%20%20%20%20L56.2%2C22.5l-0.4-0.3l-22.9%2C5l-11.7%2C2.6l-0.6%2C0.1l-0.2%2C0.2l-4.2%2C2.7l0.8-10L13%2C28.7l1.3-11.2l0-0.1L13%2C9.4L12%2C19.6L8.3%2C34.4%20%20%20%20l-8%2C16.2l2.4-1.4l3.7%2C1.9l1-4.6h0l13.4-7.8l0.1-0.1l0.3%2C0.1l13.4%2C2.2L29%2C43.6l-0.3%2C0.2l0.1%2C0.7l2.9%2C14.2l0.2%2C0.8l0.7%2C0.6l8%2C6.3%20%20%20%20l-0.2-5.7l-7.8-1.6l1-15.2l9.6-1.5l1.2-0.2l0.2%2C0.8l0.2%2C0.1l0.4%2C0.2l0.4%2C0.2l0.3%2C0.1l22.2%2C9.5L72%2C63.5l12.8%2C4.4L91%2C86.7l0.1%2C0.4%20%20%20%20l0.8%2C0.3l7.8%2C3.2L94.7%2C82.9z%20M6.7%2C45.9L6.7%2C45.9l-4.5%2C2.6l5.2-10.4L9.9%2C33l0%2C0L8.6%2C38l-0.2%2C0.7L6.7%2C45.9z%20M19.7%2C38.4l-12.1%2C7%20%20%20%20l1.7-6.9l3.6-1.3L10.7%2C33l0.7-2.9l3.7%2C3.4l0%2C0l0.2%2C0.2l0.1%2C0.1l0%2C0l0.2%2C0.2l0.5%2C0.5l0.1%2C0.1l3.8%2C3.5l0.1%2C0.1L19.7%2C38.4z%20%20%20%20%20M21.3%2C37.6v-7L34%2C27.8l7.8%2C4.5L21.3%2C37.6z%20M45.6%2C42.4l-0.3-0.1l-0.1-0.4L45%2C41.1l-0.8-4l-0.2-0.8l-0.7-3.5l8.1-0.7L45.6%2C42.4z%20%20%20%20%20M52.5%2C31.2L52.5%2C31.2L52.5%2C31.2L52.5%2C31.2z%20M52.8%2C31.6L52.8%2C31.6l0.1-0.4h0l3-7.9l12%2C10.4l-1.9%2C10.7L52.8%2C31.6z%20M74.1%2C55.8%20%20%20%20l-5.3-3.4v0l0%2C0l-2.2-7l7.5-6.2l9.4%2C8.1l0.8%2C6.1l-4.4%2C6.1L74.1%2C55.8z%20M85.8%2C68.3l4.3%2C1.5l0.8%2C14L85.8%2C68.3z%22%20%2F%3E%0A%09%09%3C%2Fg%3E%0A%09%09%3Cg%3E%0A%09%09%09%3Cpath%20d%3D%22M15.2%2C33.5L15.2%2C33.5l0.2%2C0.2L15.2%2C33.5z%22%20%2F%3E%0A%09%09%3C%2Fg%3E%0A%09%09%3Cg%3E%0A%09%09%09%3Cpath%20d%3D%22M15.2%2C33.5L15.2%2C33.5l0.2%2C0.2L15.2%2C33.5z%22%20%2F%3E%0A%09%09%3C%2Fg%3E%0A%09%09%3Cg%3E%0A%09%09%09%3Cpolygon%20points%3D%2252.5%2C31.2%2052.5%2C31.2%2052.4%2C31.2%20%20%20%22%20%2F%3E%0A%09%09%09%3Cpolygon%20points%3D%2252.5%2C31.2%2052.5%2C31.2%2052.4%2C31.2%20%20%20%22%20%2F%3E%0A%09%09%3C%2Fg%3E%0A%09%3C%2Fg%3E%0A%3C%2Fsvg%3E%0A" style="height: 140px"/>
</td>
<td><img src="/rules_proto/assets/images/protobuf.png" style="height: 180px"/></td>
<td><img src="https://avatars2.githubusercontent.com/u/7802525?v=4&s=400" style="height: 140px"/></td>
</tr><tr>
<td>Bazel</td>
<td>Gazelle</td>
<td>Protobuf</td>
<td>gRPC</td>
</tr></table>

`stackb/rules_proto` contains code and examples for using protobuf in your Bazel repository.

`@build_stack_rules_proto//language/protobuf` implements a gazelle "Language" that parses your `*.proto` files and generates BUILD rules.  You can integrate this language into your own repository via the [gazelle_binary](https://github.com/bazelbuild/bazel-gazelle/blob/master/extend.rst#gazelle_binary) rule.

The `protobuf` gazelle language is configured using gazelle "directives".
Here's an example configuration that enables a "proto_language" that combines
the "python" plugin (the one that's built-in to protoc), grpc, and
dropbox/mypy-protobuf:

```python
# -- The "proto_rule" directive instantiates a rule configuration --
# gazelle:proto_rule proto_compile implementation stackb:rules_proto:proto_compile

# -- The "proto_plugin" directive instantiates a plugin configuration --
# gazelle:proto_plugin python implementation builtin:python
# gazelle:proto_plugin grpc-python implementation grpc:grpc:grpc_python_plugin
# gazelle:proto_plugin mypy implementation dropbox:mypy-protobuf:protoc-gen-mypy

# -- The "proto_language" directive binds the rule(s) and plugin(s) together --
# gazelle:proto_language python rule proto_compile
# gazelle:proto_language python plugin python
# gazelle:proto_language python plugin grpc-python
# gazelle:proto_language python plugin mypy
# gazelle:proto_language python enabled true
```

With this configuration, any subdirectory that contains `*.proto` files will be
scanned by gazelle and corresponding proto build rules produced.

Here's what you might expect to be generated with the above configuration (NOTE:
This extension delegates to `@bazel_gazelle//language/proto` to produce the
`proto_library` rule).
:

```python
proto_library(
    name = "example_proto",
    srcs = ["example.proto"],
    visibility = ["//visibility:public"],
)

proto_compile(
    name = "example_python_compile",
    outputs = [
        "example_pb2_grpc.py",
        "example_pb2.py",
        "example_pb2.pyi",
    ],
    plugins = [
        "@build_stack_rules_proto//plugin/builtin:python",
        "@build_stack_rules_proto//plugin/grpc/grpc:grpc_python_plugin",
        "@build_stack_rules_proto//plugin/dropbox/mypy-protobuf:protoc-gen-grpc-mypy",
    ],
    proto = "example_proto",
)
```

The generated outputs could then be consumed as follows:

```python
py_library(
    srcs = [
        "example_pb2_grpc.py",
        "example_pb2.py",
    ],
    deps = [
        requirement("protobuf"),
        requirement("grpcio"),
    ],
)
```

Please refer to the following guides to get started:

- [Getting Started with Gazelle](install)
- [Getting Started without Gazelle](install)
- [Gazelle Directives Reference](config)
- [Implementing custom gazelle protobuf plugins and rules](guides/custom)
- [Contributing to rules_proto](guides/contributors)

For reference, the [examples/golden/testdata] directory contains a number of
tested example configurations.