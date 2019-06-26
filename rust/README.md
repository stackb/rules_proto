# `rust`

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
load("@build_stack_rules_proto//rust:deps.bzl", "rust_proto_compile")

rust_proto_compile()

load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")

rust_repositories()

load("@io_bazel_rules_rust//:workspace.bzl", "bazel_version")

bazel_version(name = "bazel_version")

load("@io_bazel_rules_rust//proto/raze:crates.bzl", "raze_fetch_remote_crates")

raze_fetch_remote_crates()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//rust:rust_proto_compile.bzl", "rust_proto_compile")

rust_proto_compile(
    name = "person_rust_proto",
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

## `rust_grpc_compile`

Generates rust protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//rust:deps.bzl", "rust_grpc_compile")

rust_grpc_compile()

load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")

rust_repositories()

load("@io_bazel_rules_rust//:workspace.bzl", "bazel_version")

bazel_version(name = "bazel_version")

load("@io_bazel_rules_rust//proto/raze:crates.bzl", "raze_fetch_remote_crates")

raze_fetch_remote_crates()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//rust:rust_grpc_compile.bzl", "rust_grpc_compile")

rust_grpc_compile(
    name = "greeter_rust_grpc",
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

---

## `rust_proto_library`

Generates rust protobuf library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//rust:deps.bzl", "rust_proto_library")

rust_proto_library()

load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")

rust_repositories()

load("@io_bazel_rules_rust//:workspace.bzl", "bazel_version")

bazel_version(name = "bazel_version")

load("@io_bazel_rules_rust//proto/raze:crates.bzl", "raze_fetch_remote_crates")

raze_fetch_remote_crates()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//rust:rust_proto_library.bzl", "rust_proto_library")

rust_proto_library(
    name = "person_rust_library",
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

## `rust_grpc_library`

Generates rust protobuf+gRPC library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//rust:deps.bzl", "rust_grpc_library")

rust_grpc_library()

load("@io_bazel_rules_rust//rust:repositories.bzl", "rust_repositories")

rust_repositories()

load("@io_bazel_rules_rust//:workspace.bzl", "bazel_version")

bazel_version(name = "bazel_version")

load("@io_bazel_rules_rust//proto/raze:crates.bzl", "raze_fetch_remote_crates")

raze_fetch_remote_crates()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//rust:rust_grpc_library.bzl", "rust_grpc_library")

rust_grpc_library(
    name = "greeter_rust_library",
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
