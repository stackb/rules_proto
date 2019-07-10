# Objective-C rules

| Rule | Description |
| ---: | :--- |
| [objc_proto_compile](#objc_proto_compile) | Generates objc protobuf artifacts |
| [objc_grpc_compile](#objc_grpc_compile) | Generates objc protobuf+gRPC artifacts |
| [objc_proto_library](#objc_proto_library) | Generates objc protobuf library |

---

## `objc_proto_compile`

Generates objc protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//objc:deps.bzl", "objc_deps")

objc_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//objc:defs.bzl", "objc_proto_compile")

objc_proto_compile(
    name = "person_objc_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `objc_grpc_compile`

Generates objc protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//objc:deps.bzl", "objc_deps")

objc_deps()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//objc:defs.bzl", "objc_grpc_compile")

objc_grpc_compile(
    name = "greeter_objc_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `objc_proto_library`

Generates objc protobuf library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//objc:deps.bzl", "objc_deps")

objc_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//objc:defs.bzl", "objc_proto_library")

objc_proto_library(
    name = "person_objc_library",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
