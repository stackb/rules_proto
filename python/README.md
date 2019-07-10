# Python rules

Rules for generating Python protobuf and gRPC `.py` files and libraries using standard Protocol Buffers and gRPC. Libraries are created with the Bazel native `py_library`

| Rule | Description |
| ---: | :--- |
| [python_proto_compile](#python_proto_compile) | Generates Python protobuf `.py` artifacts |
| [python_grpc_compile](#python_grpc_compile) | Generates Python protobuf+gRPC `.py` artifacts |
| [python_proto_library](#python_proto_library) | Generates a Python protobuf library using `py_library` |
| [python_grpc_library](#python_grpc_library) | Generates a Python protobuf+gRPC library using `py_library` |

---

## `python_proto_compile`

Generates Python protobuf `.py` artifacts

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

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `python_grpc_compile`

Generates Python protobuf+gRPC `.py` artifacts

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

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `python_proto_library`

Generates a Python protobuf library using `py_library`

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

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `python_grpc_library`

Generates a Python protobuf+gRPC library using `py_library`

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

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
| `python_version` | `string` | false | `PY3`    | Specify the Python version to use for the bundled dependencies. Valid values are "PY3" (the default) and "PY2"          |
