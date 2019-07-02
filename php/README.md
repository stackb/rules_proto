# `php`

| Rule | Description |
| ---: | :--- |
| [php_proto_compile](#php_proto_compile) | Generates php protobuf artifacts |
| [php_grpc_compile](#php_grpc_compile) | Generates php protobuf+gRPC artifacts |

---

## `php_proto_compile`

Generates php protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//php:deps.bzl", "php_proto_compile")

php_proto_compile()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//php:defs.bzl", "php_proto_compile")

php_proto_compile(
    name = "person_php_proto",
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

## `php_grpc_compile`

Generates php protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//php:deps.bzl", "php_grpc_compile")

php_grpc_compile()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//php:defs.bzl", "php_grpc_compile")

php_grpc_compile(
    name = "greeter_php_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
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
