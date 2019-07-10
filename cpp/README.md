# C++ rules

Rules for generating C++ protobuf and gRPC `.cc` & `.h` files and libraries using standard Protocol Buffers and gRPC. Libraries are created with the Bazel native `cc_library`

| Rule | Description |
| ---: | :--- |
| [cpp_proto_compile](#cpp_proto_compile) | Generates C++ protobuf `.h` & `.cc` artifacts |
| [cpp_grpc_compile](#cpp_grpc_compile) | Generates C++ protobuf+gRPC `.h` & `.cc` artifacts |
| [cpp_proto_library](#cpp_proto_library) | Generates a C++ protobuf library using `cc_library`, with dependencies linked |
| [cpp_grpc_library](#cpp_grpc_library) | Generates a C++ protobuf+gRPC library using `cc_library`, with dependencies linked |

---

## `cpp_proto_compile`

Generates C++ protobuf `.h` & `.cc` artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//cpp:deps.bzl", "cpp_deps")

cpp_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//cpp:defs.bzl", "cpp_proto_compile")

cpp_proto_compile(
    name = "person_cpp_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `cpp_grpc_compile`

Generates C++ protobuf+gRPC `.h` & `.cc` artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//cpp:deps.bzl", "cpp_deps")

cpp_deps()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//cpp:defs.bzl", "cpp_grpc_compile")

cpp_grpc_compile(
    name = "greeter_cpp_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `cpp_proto_library`

Generates a C++ protobuf library using `cc_library`, with dependencies linked

### `WORKSPACE`

```python
load("@build_stack_rules_proto//cpp:deps.bzl", "cpp_deps")

cpp_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//cpp:defs.bzl", "cpp_proto_library")

cpp_proto_library(
    name = "person_cpp_library",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `cpp_grpc_library`

Generates a C++ protobuf+gRPC library using `cc_library`, with dependencies linked

### `WORKSPACE`

```python
load("@build_stack_rules_proto//cpp:deps.bzl", "cpp_deps")

cpp_deps()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//cpp:defs.bzl", "cpp_grpc_library")

cpp_grpc_library(
    name = "greeter_cpp_library",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
