# Go rules

Rules for generating Go protobuf and gRPC `.go` files and libraries using [golang/protobuf](https://github.com/golang/protobuf). Libraries are created with `go_library` from [rules_go](https://github.com/bazelbuild/rules_go)

| Rule | Description |
| ---: | :--- |
| [go_proto_compile](#go_proto_compile) | Generates Go protobuf `.go` artifacts |
| [go_grpc_compile](#go_grpc_compile) | Generates Go protobuf+gRPC `.go` artifacts |
| [go_proto_library](#go_proto_library) | Generates a Go protobuf library using `go_library` from `rules_go` |
| [go_grpc_library](#go_grpc_library) | Generates a Go protobuf+gRPC library using `go_library` from `rules_go` |

---

## `go_proto_compile`

Generates Go protobuf `.go` artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//go:deps.bzl", "go_deps")

go_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//go:defs.bzl", "go_proto_compile")

go_proto_compile(
    name = "person_go_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `go_grpc_compile`

Generates Go protobuf+gRPC `.go` artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//go:deps.bzl", "go_deps")

go_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//go:defs.bzl", "go_grpc_compile")

go_grpc_compile(
    name = "greeter_go_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `go_proto_library`

Generates a Go protobuf library using `go_library` from `rules_go`

### `WORKSPACE`

```python
load("@build_stack_rules_proto//go:deps.bzl", "go_deps")

go_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//go:defs.bzl", "go_proto_library")

go_proto_library(
    name = "person_go_library",
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/stackb/rules_proto/go/example/go_proto_library/person",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| `importpath` | `string` | false | `None`    | Importpath for the generated artifacts          |

---

## `go_grpc_library`

Generates a Go protobuf+gRPC library using `go_library` from `rules_go`

### `WORKSPACE`

```python
load("@build_stack_rules_proto//go:deps.bzl", "go_deps")

go_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//go:defs.bzl", "go_grpc_library")

go_grpc_library(
    name = "greeter_go_library",
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/stackb/rules_proto/go/example/go_grpc_library/greeter",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| `importpath` | `string` | false | `None`    | Importpath for the generated artifacts          |
