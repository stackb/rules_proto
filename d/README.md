# `d`

These rules use the protoc-gen-d plugin, which only supports proto3 .proto files.

| Rule | Description |
| ---: | :--- |
| [d_proto_compile](#d_proto_compile) | Generates d protobuf artifacts |
| [d_proto_library](#d_proto_library) | Generates d protobuf library |

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
load("@build_stack_rules_proto//d:defs.bzl", "d_proto_compile")

d_proto_compile(
    name = "person_d_proto",
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
| verbose   | `int` | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `d_proto_library`

Generates d protobuf library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//d:deps.bzl", "d_proto_library")

d_proto_library()

load("@io_bazel_rules_d//d:d.bzl", "d_repositories")

d_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//d:defs.bzl", "d_proto_library")

d_proto_library(
    name = "person_d_library",
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
| verbose   | `int` | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
