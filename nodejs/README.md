# `nodejs`

| Rule | Description |
| ---: | :--- |
| [nodejs_proto_compile](#nodejs_proto_compile) | Generates node *.js protobuf artifacts |
| [nodejs_grpc_compile](#nodejs_grpc_compile) | Generates node *.js protobuf+gRPC artifacts |

---

## `nodejs_proto_compile`

Generates node *.js protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//nodejs:deps.bzl", "nodejs_deps")

nodejs_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//nodejs:defs.bzl", "nodejs_proto_compile")

nodejs_proto_compile(
    name = "person_nodejs_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `nodejs_grpc_compile`

Generates node *.js protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//nodejs:deps.bzl", "nodejs_deps")

nodejs_deps()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//nodejs:defs.bzl", "nodejs_grpc_compile")

nodejs_grpc_compile(
    name = "greeter_nodejs_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
