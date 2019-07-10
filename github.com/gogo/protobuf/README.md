# Go (gogoprotobuf) rules

Rules for generating Go protobuf and gRPC `.go` files and libraries using [gogo/protobuf](https://github.com/gogo/protobuf). Libraries are created with `go_library` from [rules_go](https://github.com/bazelbuild/rules_go)

| Rule | Description |
| ---: | :--- |
| [gogo_proto_compile](#gogo_proto_compile) | Generates gogo protobuf `.go` artifacts |
| [gogo_grpc_compile](#gogo_grpc_compile) | Generates gogo protobuf+gRPC `.go` artifacts |
| [gogo_proto_library](#gogo_proto_library) | Generates a Go gogo protobuf library using `go_library` from `rules_go` |
| [gogo_grpc_library](#gogo_grpc_library) | Generates a Go gogo protobuf+gRPC library using `go_library` from `rules_go` |
| [gogofast_proto_compile](#gogofast_proto_compile) | Generates gogofast protobuf `.go` artifacts |
| [gogofast_grpc_compile](#gogofast_grpc_compile) | Generates gogofast protobuf+gRPC `.go` artifacts |
| [gogofast_proto_library](#gogofast_proto_library) | Generates a Go gogofast protobuf library using `go_library` from `rules_go` |
| [gogofast_grpc_library](#gogofast_grpc_library) | Generates a Go gogofast protobuf+gRPC library using `go_library` from `rules_go` |
| [gogofaster_proto_compile](#gogofaster_proto_compile) | Generates gogofaster protobuf `.go` artifacts |
| [gogofaster_grpc_compile](#gogofaster_grpc_compile) | Generates gogofaster protobuf+gRPC `.go` artifacts |
| [gogofaster_proto_library](#gogofaster_proto_library) | Generates a Go gogofaster protobuf library using `go_library` from `rules_go` |
| [gogofaster_grpc_library](#gogofaster_grpc_library) | Generates a Go gogofaster protobuf+gRPC library using `go_library` from `rules_go` |

---

## `gogo_proto_compile`

Generates gogo protobuf `.go` artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_deps")

gogo_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:defs.bzl", "gogo_proto_compile")

gogo_proto_compile(
    name = "person_gogo_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `gogo_grpc_compile`

Generates gogo protobuf+gRPC `.go` artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_deps")

gogo_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:defs.bzl", "gogo_grpc_compile")

gogo_grpc_compile(
    name = "greeter_gogo_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `gogo_proto_library`

Generates a Go gogo protobuf library using `go_library` from `rules_go`

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_deps")

gogo_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:defs.bzl", "gogo_proto_library")

gogo_proto_library(
    name = "person_gogo_library",
    go_deps = [
        "@com_github_gogo_protobuf//types:go_default_library",
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/stackb/rules_proto/gogo/example/gogo_proto_library/person",
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

## `gogo_grpc_library`

Generates a Go gogo protobuf+gRPC library using `go_library` from `rules_go`

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_deps")

gogo_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:defs.bzl", "gogo_grpc_library")

gogo_grpc_library(
    name = "greeter_gogo_library",
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/stackb/rules_proto/gogo/example/gogo_grpc_library/greeter",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| `importpath` | `string` | false | `None`    | Importpath for the generated artifacts          |

---

## `gogofast_proto_compile`

Generates gogofast protobuf `.go` artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_deps")

gogo_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:defs.bzl", "gogofast_proto_compile")

gogofast_proto_compile(
    name = "person_gogo_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `gogofast_grpc_compile`

Generates gogofast protobuf+gRPC `.go` artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_deps")

gogo_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:defs.bzl", "gogofast_grpc_compile")

gogofast_grpc_compile(
    name = "greeter_gogo_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `gogofast_proto_library`

Generates a Go gogofast protobuf library using `go_library` from `rules_go`

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_deps")

gogo_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:defs.bzl", "gogofast_proto_library")

gogofast_proto_library(
    name = "person_gogo_library",
    go_deps = [
        "@com_github_gogo_protobuf//types:go_default_library",
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/stackb/rules_proto/gogo/example/gogofast_proto_library/person",
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

## `gogofast_grpc_library`

Generates a Go gogofast protobuf+gRPC library using `go_library` from `rules_go`

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_deps")

gogo_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:defs.bzl", "gogofast_grpc_library")

gogofast_grpc_library(
    name = "greeter_gogo_library",
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/stackb/rules_proto/gogo/example/gogofast_grpc_library/greeter",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| `importpath` | `string` | false | `None`    | Importpath for the generated artifacts          |

---

## `gogofaster_proto_compile`

Generates gogofaster protobuf `.go` artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_deps")

gogo_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:defs.bzl", "gogofaster_proto_compile")

gogofaster_proto_compile(
    name = "person_gogo_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `gogofaster_grpc_compile`

Generates gogofaster protobuf+gRPC `.go` artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_deps")

gogo_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:defs.bzl", "gogofaster_grpc_compile")

gogofaster_grpc_compile(
    name = "greeter_gogo_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `gogofaster_proto_library`

Generates a Go gogofaster protobuf library using `go_library` from `rules_go`

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_deps")

gogo_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:defs.bzl", "gogofaster_proto_library")

gogofaster_proto_library(
    name = "person_gogo_library",
    go_deps = [
        "@com_github_gogo_protobuf//types:go_default_library",
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/stackb/rules_proto/gogo/example/gogofaster_proto_library/person",
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

## `gogofaster_grpc_library`

Generates a Go gogofaster protobuf+gRPC library using `go_library` from `rules_go`

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:deps.bzl", "gogo_deps")

gogo_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/gogo/protobuf:defs.bzl", "gogofaster_grpc_library")

gogofaster_grpc_library(
    name = "greeter_gogo_library",
    go_deps = [
        "@com_github_golang_protobuf//ptypes/any:go_default_library",
    ],
    importpath = "github.com/stackb/rules_proto/gogo/example/gogofaster_grpc_library/greeter",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| `importpath` | `string` | false | `None`    | Importpath for the generated artifacts          |
