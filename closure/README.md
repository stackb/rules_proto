# `closure`

| Rule | Description |
| ---: | :--- |
| [closure_proto_compile](#closure_proto_compile) | Generates closure *.js protobuf+gRPC files |
| [closure_proto_library](#closure_proto_library) | Generates a closure_library with compiled protobuf *.js files |

---

## `closure_proto_compile`

Generates closure *.js protobuf+gRPC files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//closure:deps.bzl", "closure_proto_compile")

closure_proto_compile()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//closure:closure_proto_compile.bzl", "closure_proto_compile")

closure_proto_compile(
    name = "person_closure_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Mandatory Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| deps   | `list<ProtoInfo>` | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |

### Optional Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| plugins   | `list<ProtoPluginInfo>` | `[]`    | List of labels that provide a `ProtoPluginInfo`          |
| plugin_options   | `list<string>` | `[]`    | List of additional 'global' plugin options (applies to all plugins). To apply plugin specific options, use the `options` attribute on `proto_plugin`          |
| outputs   | `list<generated file>` | `[]`    | List of additional expected generated file outputs          |
| verbose   | `int` | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |

---

## `closure_proto_library`

Generates a closure_library with compiled protobuf *.js files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//closure:deps.bzl", "closure_proto_library")

closure_proto_library()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//closure:closure_proto_library.bzl", "closure_proto_library")

closure_proto_library(
    name = "person_closure_library",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Mandatory Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| deps   | `list<ProtoInfo>` | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |

### Optional Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| plugins   | `list<ProtoPluginInfo>` | `[]`    | List of labels that provide a `ProtoPluginInfo`          |
| plugin_options   | `list<string>` | `[]`    | List of additional 'global' plugin options (applies to all plugins). To apply plugin specific options, use the `options` attribute on `proto_plugin`          |
| outputs   | `list<generated file>` | `[]`    | List of additional expected generated file outputs          |
| verbose   | `int` | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
