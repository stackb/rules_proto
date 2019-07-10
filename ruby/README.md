# Ruby rules

Rules for generating Ruby protobuf and gRPC `.rb` files and libraries using standard Protocol Buffers and gRPC. Libraries are created with `ruby_library` from [rules_ruby](https://github.com/yugui/rules_ruby)

| Rule | Description |
| ---: | :--- |
| [ruby_proto_compile](#ruby_proto_compile) | Generates Ruby protobuf `.rb` artifacts |
| [ruby_grpc_compile](#ruby_grpc_compile) | Generates Ruby protobuf+gRPC `.rb` artifacts |
| [ruby_proto_library](#ruby_proto_library) | Generates a Ruby protobuf library using `ruby_library` from `rules_ruby` |
| [ruby_grpc_library](#ruby_grpc_library) | Generates a Ruby protobuf+gRPC library using `ruby_library` from `rules_ruby` |

---

## `ruby_proto_compile`

Generates Ruby protobuf `.rb` artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//ruby:deps.bzl", "ruby_deps")

ruby_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//ruby:defs.bzl", "ruby_proto_compile")

ruby_proto_compile(
    name = "person_ruby_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `ruby_grpc_compile`

Generates Ruby protobuf+gRPC `.rb` artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//ruby:deps.bzl", "ruby_deps")

ruby_deps()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//ruby:defs.bzl", "ruby_grpc_compile")

ruby_grpc_compile(
    name = "greeter_ruby_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `ruby_proto_library`

Generates a Ruby protobuf library using `ruby_library` from `rules_ruby`

### `WORKSPACE`

```python
load("@build_stack_rules_proto//ruby:deps.bzl", "ruby_deps")

ruby_deps()

load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_register_toolchains")

ruby_register_toolchains()

load("@com_github_yugui_rules_ruby//ruby/private:bundle.bzl", "bundle_install")

bundle_install(
    name = "routeguide_gems_bundle",
    gemfile = "//ruby:Gemfile",
    gemfile_lock = "//ruby:Gemfile.lock",
)
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//ruby:defs.bzl", "ruby_proto_library")

ruby_proto_library(
    name = "person_ruby_library",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `ruby_grpc_library`

Generates a Ruby protobuf+gRPC library using `ruby_library` from `rules_ruby`

### `WORKSPACE`

```python
load("@build_stack_rules_proto//ruby:deps.bzl", "ruby_deps")

ruby_deps()

load("@com_github_yugui_rules_ruby//ruby:def.bzl", "ruby_register_toolchains")

ruby_register_toolchains()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()

load("@com_github_yugui_rules_ruby//ruby/private:bundle.bzl", "bundle_install")

bundle_install(
    name = "routeguide_gems_bundle",
    gemfile = "//ruby:Gemfile",
    gemfile_lock = "//ruby:Gemfile.lock",
)
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//ruby:defs.bzl", "ruby_grpc_library")

ruby_grpc_library(
    name = "greeter_ruby_library",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
