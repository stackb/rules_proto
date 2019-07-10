# Android rules

Rules for generating Android protobuf and gRPC `.jar` files and libraries using standard Protocol Buffers and [gRPC-Java](https://github.com/grpc/grpc-java). Libraries are created with `android_library` from [rules_android](https://github.com/bazelbuild/rules_android)

| Rule | Description |
| ---: | :--- |
| [android_proto_compile](#android_proto_compile) | Generates an Android protobuf `.jar` artifact |
| [android_grpc_compile](#android_grpc_compile) | Generates Android protobuf+gRPC `.jar` artifacts |
| [android_proto_library](#android_proto_library) | Generates an Android protobuf library using `android_library` from `rules_android` |
| [android_grpc_library](#android_grpc_library) | Generates Android protobuf+gRPC library using `android_library` from `rules_android` |

---

## `android_proto_compile`

Generates an Android protobuf `.jar` artifact

### `WORKSPACE`

```python
load("@build_stack_rules_proto//android:deps.bzl", "android_deps")

android_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//android:defs.bzl", "android_proto_compile")

android_proto_compile(
    name = "person_android_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `android_grpc_compile`

Generates Android protobuf+gRPC `.jar` artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//android:deps.bzl", "android_deps")

android_deps()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories(
    omit_com_google_protobuf = True,
    omit_net_zlib = True
)
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//android:defs.bzl", "android_grpc_compile")

android_grpc_compile(
    name = "greeter_android_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `android_proto_library`

Generates an Android protobuf library using `android_library` from `rules_android`

### `WORKSPACE`

```python
# The set of dependencies loaded here is excessive for android proto alone
# (but simplifies our setup)
load("@build_stack_rules_proto//android:deps.bzl", "android_deps")

android_deps()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories(
    omit_com_google_protobuf = True,
    omit_com_google_protobuf_javalite = True,
    omit_net_zlib = True
)

load("@build_bazel_rules_android//android:sdk_repository.bzl", "android_sdk_repository")

android_sdk_repository(name = "androidsdk")
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//android:defs.bzl", "android_proto_library")

android_proto_library(
    name = "person_android_library",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `android_grpc_library`

Generates Android protobuf+gRPC library using `android_library` from `rules_android`

### `WORKSPACE`

```python
load("@build_stack_rules_proto//android:deps.bzl", "android_deps")

android_deps()

load("@io_grpc_grpc_java//:repositories.bzl", "grpc_java_repositories")

grpc_java_repositories(
    omit_com_google_protobuf = True,
    omit_com_google_protobuf_javalite = True,
    omit_net_zlib = True
)

load("@build_bazel_rules_android//android:sdk_repository.bzl", "android_sdk_repository")

android_sdk_repository(name = "androidsdk")
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//android:defs.bzl", "android_grpc_library")

android_grpc_library(
    name = "greeter_android_library",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
