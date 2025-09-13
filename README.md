# `rules_proto (v3)`

![build-status](https://github.com/stackb/rules_proto/actions/workflows/ci.yaml/badge.svg)
[![Go Reference](https://pkg.go.dev/badge/github.com/stackb/rules_proto.svg)](https://pkg.go.dev/github.com/stackb/rules_proto)

Bazel starlark rules for building protocol buffers +/- gRPC :sparkles:.

<table border="0">
  <tr>
    <td><img src="https://upload.wikimedia.org/wikipedia/commons/7/7d/Bazel_logo.svg" height="120"/></td>
    <td><img src="https://user-images.githubusercontent.com/50580/141892423-5205bbfd-8487-442b-81c7-f56fa3d1f69e.jpeg" height="120"/></td>
    <td><img src="https://user-images.githubusercontent.com/50580/141900696-bfb2d42d-5d2c-46f8-bd9f-06515969f6a2.png" height="120"/></td>
    <td><img src="https://avatars2.githubusercontent.com/u/7802525?v=4&s=400" height="120"/></td>
  </tr>
  <tr>
    <td>bazel</td>
    <td>gazelle</td>
    <td>protobuf</td>
    <td>grpc</td>
  </tr>
</table>

`@build_stack_rules_proto` provides:

1. Rules for driving the `protoc` tool within a bazel workspace.
2. A [gazelle](https://github.com/bazelbuild/bazel-gazelle/) extension that
   generates rules based on the content of your `.proto` files.
3. Example setups for a variety of languages.

## `MODULE.bazel`

```py
bazel_dep(name = "build_stack_rules_proto", version = "3.0.0")
```

See <https://registry.bazel.build/search?q=rules_proto> for latest version.

> NOTE: Version 3.x.x no longer supports `WORKSPACE`, please use the latest
> 2.x.x release for workspace compatibility.

## Docs

| Description                               | Link                                   |
|-------------------------------------------|----------------------------------------|
| For documentation about the core ruleset  | [BUILD_RULES.md](/docs/BUILD_RULES.md) |
| For info about toolchains                 | [TOOLCHAINS.md](/docs/TOOLCHAINS.md)   |
| For help setting up the gazelle extension | [GAZELLE.md](/docs/GAZELLE.md)         |
| Writing custom gazelle logic in starlark  | [STARLARK.md](/docs/STARLARK.md)       |
| For a history of this repository          | [HISTORY.md](/docs/HISTORY.md)         |

## Examples

> Examples are generated from golden tests

| Description      | Link                                          |
|------------------|-----------------------------------------------|
| commonjs         | [commonjs](/docs/commonjs.md)                 |
| cpp              | [cpp](/docs/cpp.md)                           |
| csharp           | [csharp](/docs/csharp.md)                     |
| go               | [go](/docs/go.md)                             |
| java             | [java](/docs/java.md)                         |
| objc             | [objc](/docs/objc.md)                         |
| ruby             | [ruby](/docs/ruby.md)                         |
| python           | [python](/docs/python.md)                     |
| scala            | [scala](/docs/scala.md)                       |
| proto_repository | [proto_repository](/docs/proto_repository.md) |
| starlark_java    | [starlark_java](/docs/starlark_java.md)       |
