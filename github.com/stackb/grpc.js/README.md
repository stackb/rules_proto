# `grpc.js`

| Rule | Description |
| ---: | :--- |
| [closure_grpc_compile](#closure_grpc_compile) | Generates protobuf closure grpc *.js files |
| [closure_grpc_library](#closure_grpc_library) | Generates protobuf closure library *.js files |

---

## `closure_grpc_compile`

Generates protobuf closure grpc *.js files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/stackb/grpc.js:deps.bzl", "closure_grpc_compile")

closure_grpc_compile()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/stackb/grpc.js:closure_grpc_compile.bzl", "closure_grpc_compile")

closure_grpc_compile(
    name = "greeter_grpc.js_grpc",
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

## `closure_grpc_library`

Generates protobuf closure library *.js files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/stackb/grpc.js:deps.bzl", "closure_grpc_library")

closure_grpc_library()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/stackb/grpc.js:closure_grpc_library.bzl", "closure_grpc_library")

closure_grpc_library(
    name = "greeter_grpc.js_library",
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
