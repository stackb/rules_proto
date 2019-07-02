# `node`

| Rule | Description |
| ---: | :--- |
| [node_proto_compile](#node_proto_compile) | Generates node *.js protobuf artifacts |
| [node_grpc_compile](#node_grpc_compile) | Generates node *.js protobuf+gRPC artifacts |

---

## `node_proto_compile`

Generates node *.js protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//node:deps.bzl", "node_proto_compile")

node_proto_compile()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//node:node_proto_compile.bzl", "node_proto_compile")

node_proto_compile(
    name = "person_node_proto",
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

## `node_grpc_compile`

Generates node *.js protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//node:deps.bzl", "node_grpc_compile")

node_grpc_compile()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//node:node_grpc_compile.bzl", "node_grpc_compile")

node_grpc_compile(
    name = "greeter_node_grpc",
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
