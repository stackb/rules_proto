# `rules_proto (v4)`

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
bazel_dep(name = "build_stack_rules_proto", version = "4.x.x")
```

> NOTE: `build_stack_rules_proto` is still in the submission process to the bcr.
> Until merged, a git_override or other override is needed to consume this
> repository.

See <https://registry.bazel.build/modules/build_stack_rules_proto> for latest version.

> NOTE: Version 4.x.x no longer supports `WORKSPACE`, please use the latest
> 3.x.x release for workspace compatibility.

## Docs

| Description                               | Link                                    |
|-------------------------------------------|-----------------------------------------|
| Documentation about the core ruleset      | [CORE_RULES.md](/docs/CORE_RULES.md)    |
| Available Toolchains                      | [TOOLCHAINS.md](/docs/TOOLCHAINS.md)    |
| Guide to setting up the gazelle extension | [GAZELLE.md](/docs/GAZELLE.md)          |
| Writing custom gazelle logic              | [STARLARK.md](/docs/STARLARK.md)        |
| Examples                                  | [example/README.md](/example/README.md) |
| Preconfigured plugins                     | [PLUGINS.md](/docs/PLUGINS.md)          |
| Preconfigured rules                       | [RULES.md](/docs/RULES.md)              |
| A history of this repository              | [HISTORY.md](/docs/HISTORY.md)          |
| Developer Guide                           | [DEVELOPMENT.md](/docs/DEVELOPMENT.md)  |
