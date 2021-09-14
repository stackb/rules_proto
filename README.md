# `rules_proto`

Bazel starlark rules for building protocol buffers +/- gRPC :sparkles:.

<table border="0"><tr>
<td><img src="https://bazel.build/images/bazel-icon.svg" height="180"/></td>
<td><img src="https://github.com/pubref/rules_protobuf/blob/master/images/wtfcat.png" height="180"/></td>
<td><img src="https://avatars2.githubusercontent.com/u/7802525?v=4&s=400" height="180"/></td>
</tr><tr>
<td>Bazel</td>
<td>rules_proto</td>
<td>gRPC</td>
</tr></table>

`stackb/rules_proto` combines [Bazel](https://bazel.build) rules for building
protobuf and gRPC outputs with
[gazelle](https://github.com/bazelbuild/bazel-gazelle) extension that
autogenerates those rules for you :smile:.

## Core Rules

| Rule                     | Description                                                             |
|--------------------------|-------------------------------------------------------------------------|
| `proto_plugin`           | Provides `proto_compile` with plugin-specific configuration.            |
| `proto_compile`          | bazel rule that drives the `protoc` tool                                |
| `proto_compiled_sources` | runs `proto_compile`; copies generated sources back into the workspace. |

## Language-Specific Rules

| Rule                    | Description                                 |
|-------------------------|---------------------------------------------|
| `proto_cc_library`      | protobuf-specific wrapper for `cc_library`. |
| `proto_grpc_cc_library` | gRPC-specific wrapper for `cc_library`.     |

> with the capability of copying generated sources back into the workspace (for those use-cases when you need to commit generated files under source control

## Guides

| Guide               | Description                                   |
|---------------------|-----------------------------------------------|
| [Getting Started]() | How to setup your `WORKSPACE` for rules_proto |

