# `python`

| Rule | Description |
| ---: | :--- |
| [python_proto_compile](#python_proto_compile) | Generates *.py protobuf artifacts |
| [python_grpc_compile](#python_grpc_compile) | Generates *.py protobuf+gRPC artifacts |
| [python_proto_library](#python_proto_library) | Generates *.py protobuf library |
| [python_grpc_library](#python_grpc_library) | Generates *.py protobuf+gRPC library |

---

## `python_proto_compile`

Generates *.py protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//python:deps.bzl", "python_deps")

python_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//python:defs.bzl", "python_proto_compile")

python_proto_compile(
    name = "person_python_proto",
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

## `python_grpc_compile`

Generates *.py protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//python:deps.bzl", "python_deps")

python_deps()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//python:defs.bzl", "python_grpc_compile")

python_grpc_compile(
    name = "greeter_python_grpc",
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

## `python_proto_library`

Generates *.py protobuf library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//python:deps.bzl", "python_deps")

python_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//python:defs.bzl", "python_proto_library")

python_proto_library(
    name = "person_python_library",
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

## `python_grpc_library`

Generates *.py protobuf+gRPC library

### `WORKSPACE`

```python
load("@build_stack_rules_proto//python:deps.bzl", "python_deps")

python_deps()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@com_apt_itude_rules_pip//rules:dependencies.bzl", "pip_rules_dependencies")

pip_rules_dependencies()

load("@com_apt_itude_rules_pip//rules:repository.bzl", "pip_repository")

pip_repository(
    name = "grpc_py2_deps",
    python_interpreter = "python2",
    requirements = "@build_stack_rules_proto//python/requirements:grpc.txt",
)

pip_repository(
    name = "grpc_py3_deps",
    python_interpreter = "python3",
    requirements = "@build_stack_rules_proto//python/requirements:grpc.txt",
)
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//python:defs.bzl", "python_grpc_library")

python_grpc_library(
    name = "greeter_python_library",
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
| python_version   | `string` | `PY3`    | Specify the Python version to use for the bundled dependencies. Valid values are "PY3" (the default) and "PY2"          |
