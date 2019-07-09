# `grpc.js`

| Rule | Description |
| ---: | :--- |
| [grpcjs_grpc_compile](#grpcjs_grpc_compile) | Generates protobuf closure grpc *.js files |
| [grpcjs_grpc_library](#grpcjs_grpc_library) | Generates protobuf closure library *.js files |

---

## `grpcjs_grpc_compile`

Generates protobuf closure grpc *.js files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/stackb/grpc.js:deps.bzl", "grpcjs_deps")

grpcjs_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/stackb/grpc.js:defs.bzl", "grpcjs_grpc_compile")

grpcjs_grpc_compile(
    name = "greeter_grpc.js_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `grpcjs_grpc_library`

Generates protobuf closure library *.js files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/stackb/grpc.js:deps.bzl", "grpcjs_deps")

grpcjs_deps()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/stackb/grpc.js:defs.bzl", "grpcjs_grpc_library")

grpcjs_grpc_library(
    name = "greeter_grpc.js_library",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
