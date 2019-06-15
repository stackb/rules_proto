# `d`

| Rule | Description |
| ---: | :--- |
| [d_proto_compile](#d_proto_compile) | Generates d protobuf artifacts |

---

## `d_proto_compile`

Generates d protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//d:deps.bzl", "d_proto_compile")

d_proto_compile()

load("@io_bazel_rules_d//d:d.bzl", "d_repositories")

d_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//d:d_proto_compile.bzl", "d_proto_compile")

d_proto_compile(
    name = "person_d_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Mandatory Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| deps   | `list<ProtoInfo>` | `[]`    | List of labels that provide a `ProtoInfo` (`native.proto_library`)          |

### Optional Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| plugins   | `list<ProtoPluginInfo>` | `[]`    | List of labels that provide a `ProtoPluginInfo`          |
| plugin_options   | `list<string>` | `[]`    | List of additional 'global' plugin options (applies to all plugins)          |
| outputs   | `list<generated file>` | `[]`    | List of additional expected generated file outputs          |
| protoc   | `executable file` | `@com_google_protobuf//:protoc`    | The protocol compiler tool          |
| verbose   | `int` | `0`    | 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
| transitive   | `bool` | `True`    | Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies          |
