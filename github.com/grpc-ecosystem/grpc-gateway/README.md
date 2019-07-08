# `grpc-gateway`

| Rule | Description |
| ---: | :--- |
| [gateway_grpc_compile](#gateway_grpc_compile) | Generates grpc-gateway *.go files |
| [gateway_swagger_compile](#gateway_swagger_compile) | Generates grpc-gateway swagger *.json files |
| [gateway_grpc_library](#gateway_grpc_library) | Generates grpc-gateway library files |

---

## `gateway_grpc_compile`

Generates grpc-gateway *.go files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//:deps.bzl", "bazel_gazelle", "io_bazel_rules_go")

io_bazel_rules_go()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

bazel_gazelle()

load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:deps.bzl", "gateway_deps")

gateway_deps()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:defs.bzl", "gateway_grpc_compile")

gateway_grpc_compile(
    name = "api_gateway_grpc",
    deps = ["@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway/example/api:api_proto"],
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

## `gateway_swagger_compile`

Generates grpc-gateway swagger *.json files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//:deps.bzl", "bazel_gazelle", "io_bazel_rules_go")

io_bazel_rules_go()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

bazel_gazelle()

load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:deps.bzl", "gateway_deps")

gateway_deps()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:defs.bzl", "gateway_swagger_compile")

gateway_swagger_compile(
    name = "api_gateway_grpc",
    deps = ["@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway/example/api:api_proto"],
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

## `gateway_grpc_library`

Generates grpc-gateway library files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//:deps.bzl", "bazel_gazelle", "io_bazel_rules_go")

io_bazel_rules_go()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains()

bazel_gazelle()

load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:deps.bzl", "gateway_deps")

gateway_deps()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway:defs.bzl", "gateway_grpc_library")

gateway_grpc_library(
    name = "api_gateway_library",
    importpath = "github.com/stackb/rules_proto/github.com/grpc-ecosystem/grpc-gateway/examples/api",
    deps = ["@build_stack_rules_proto//github.com/grpc-ecosystem/grpc-gateway/example/api:api_proto"],
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
