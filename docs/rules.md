---
layout: default
title: Rules
permalink: rules
has_children: true
nav_order: 3
---

# Rules

## `proto_compile`

Coordinates the exection of `protoc` and associated plugins.

Example:

```python
proto_compile(
    name = "example_java_compile",
    outs = {"@build_stack_rules_proto//plugin/builtin:java": "example.srcjar"},
    outputs = ["example.srcjar"],
    plugins = ["@build_stack_rules_proto//plugin/builtin:java"],
    proto = "example_proto",
)
```

| proto_compile attribute | type                      | description                                                                                         |
| ----------------------- | ------------------------- | --------------------------------------------------------------------------------------------------- |
| `plugins`               | `label_list`              | list of upstream `proto_plugin` rules (`ProtoPluginInfo` provider)                                  |
| `option`                | `label_keyed_string_list` | List of options, grouped by plugin                                                                  |
| `outs`                  | `label_keyed_string_list` | Output location for generated files, grouped by plugin. By default, the bazel execroot used.        |
| `outputs`               | `output_list`             | List of files that to be generated.                                                                 |
| `mappings`              | `label_keyed_string_list` | Mapping from execroot-relative path to package-relative path (^1).                                  |
| `srcs`                  | `string_list`             | List of source files expected to be generated. Only one of `srcs` or `outputs` should be specified. |
| `proto`                 | `label`                   | upstream `proto_library` rule (`ProtoInfo` provider)                                                |
| `verbose`               | `bool`                    | If true, the full command args, expected outputs, and pre-/post- state of sandbox are printed.      |

^1: bazel mandates that any file produced by an action is created within its
respective package path of the sandbox.  For example, consider a file
`proto/example.proto` that contains a `go_package` option `github.com/foo/bar`.
In this case the generated execroot-relative location will be
`./github.com/foo/bar/example.pb.go`, which is not inside `proto/`.  The mapping
option provides a mechanism to schedule a file copy operation `cp
/github.com/foo/bar/example.pb.go ./proto/example.pb.go` in order to satify
bazel action constraints.

## `proto_compiled_sources`

The `proto_compiled_sources` rule is intended for the use case where generated
files are checked into source control.  While one can debate whether having
generated files in the source tree is a bad idea, these files may need to remain
in source control during a migration period or because reasons.

Example:

```python
proto_compiled_sources(
    name = "v1_gogofast_compiled_sources",
    srcs = [
        "message.pb.go",
        "service.pb.go",
    ],
    output_mappings = [
        "message.pb.go=github.com/example/repo/api/v1/message.pb.go",
        "service.pb.go=github.com/example/repo/api/v1/service.pb.go",
    ],
    options = {"@build_stack_rules_proto//plugin/gogo/protobuf:protoc-gen-gogofast": ["plugins=grpc"]},
    plugins = ["@build_stack_rules_proto//plugin/gogo/protobuf:protoc-gen-gogofast"],
    proto = "v1_proto",
    visibility = ["//proto:__subpackages__"],
)
```

A `proto_compiled_sources` rule named
`//proto/api/v1:v1_gogofast_compiled_sources` is a macro that has three targets:

- `bazel build //proto/api/v1:v1_gogofast_compiled_sources` generates files in
  the `bazel-out/` directory.
- `bazel run //proto/api/v1:v1_gogofast_compiled_sources.update` generates files
  in the `bazel-out/` directory and copies them back into the source package.
- `bazel test //proto/api/v1:v1_gogofast_compiled_sources_test` asserts equality
  between the source file(s) and the generated file(s).

The `_test` target ensures that changes made to source `.proto` file will not
pass CI unless the `.update` target has been run, preventing drift.

This rule has nearly identical attributes as `proto_compile`, but the `srcs`
attribute is used rather than `outputs`.