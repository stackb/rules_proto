# Rust rules

Rules for generating Rust protobuf and gRPC `.rs` files and libraries using [rust-protobuf](https://github.com/stepancheg/rust-protobuf) and [grpc-rs](https://github.com/pingcap/grpc-rs). Libraries are created with `rust_library` from [rules_rust](https://github.com/bazelbuild/rules_rust). Due to upstream requirements, these rules require that the system has a valid GOPATH set.

| Rule | Description |
| ---: | :--- |
| [rust_proto_compile](#rust_proto_compile) | Generates Rust protobuf `.rs` artifacts |
| [rust_grpc_compile](#rust_grpc_compile) | Generates Rust protobuf+gRPC `.rs` artifacts |
| [rust_proto_library](#rust_proto_library) | Generates a Rust protobuf library using `rust_library` from `rules_rust` |
| [rust_grpc_library](#rust_grpc_library) | Generates a Rust protobuf+gRPC library using `rust_library` from `rules_rust` |

---

## `rust_proto_compile`

Generates Rust protobuf `.rs` artifacts

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

Generates Rust protobuf+gRPC `.rs` artifacts

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

Generates a Rust protobuf library using `rust_library` from `rules_rust`

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

Generates a Rust protobuf+gRPC library using `rust_library` from `rules_rust`

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
