# Rust rules

| Rule | Description |
| ---: | :--- |
| [rust_proto_compile](#rust_proto_compile) | Generates rust protobuf artifacts |
| [rust_grpc_compile](#rust_grpc_compile) | Generates rust protobuf+gRPC artifacts |
| [rust_proto_library](#rust_proto_library) | Generates rust protobuf library |
| [rust_grpc_library](#rust_grpc_library) | Generates rust protobuf+gRPC library |

---

## `rust_proto_compile`

Generates rust protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//rust:deps.bzl", "rust_deps")

rust_deps()

load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")

rust_repositories()

load("@io_bazel_rules_rust//:workspace.bzl", "bazel_version")

bazel_version(name = "bazel_version")

load("@io_bazel_rules_rust//proto:repositories.bzl", "rust_proto_repositories")

rust_proto_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//rust:defs.bzl", "rust_proto_compile")

rust_proto_compile(
    name = "person_rust_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `rust_grpc_compile`

Generates rust protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//rust:deps.bzl", "rust_deps")

rust_deps()

load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")

rust_repositories()

load("@io_bazel_rules_rust//:workspace.bzl", "bazel_version")

bazel_version(name = "bazel_version")

load("@io_bazel_rules_rust//proto:repositories.bzl", "rust_proto_repositories")

rust_proto_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//rust:defs.bzl", "rust_grpc_compile")

rust_grpc_compile(
    name = "greeter_rust_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `rust_proto_library`

Generates rust protobuf library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//rust:deps.bzl", "rust_deps")

rust_deps()

load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")

rust_repositories()

load("@io_bazel_rules_rust//:workspace.bzl", "bazel_version")

bazel_version(name = "bazel_version")

load("@io_bazel_rules_rust//proto:repositories.bzl", "rust_proto_repositories")

rust_proto_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//rust:defs.bzl", "rust_proto_library")

rust_proto_library(
    name = "person_rust_library",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `rust_grpc_library`

Generates rust protobuf+gRPC library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//rust:deps.bzl", "rust_deps")

rust_deps()

load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")

rust_repositories()

load("@io_bazel_rules_rust//:workspace.bzl", "bazel_version")

bazel_version(name = "bazel_version")

load("@io_bazel_rules_rust//proto:repositories.bzl", "rust_proto_repositories")

rust_proto_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//rust:defs.bzl", "rust_grpc_library")

rust_grpc_library(
    name = "greeter_rust_library",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
