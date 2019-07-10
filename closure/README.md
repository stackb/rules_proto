# Closure rules

Rules for generating Closure protobuf `.js` files and libraries using standard Protocol Buffers. Libraries are created with `closure_js_library` from [rules_closure](https://github.com/bazelbuild/rules_closure)

| Rule | Description |
| ---: | :--- |
| [closure_proto_compile](#closure_proto_compile) | Generates Closure protobuf `.js` files |
| [closure_proto_library](#closure_proto_library) | Generates a Closure library with compiled protobuf `.js` files using `closure_js_library` from `rules_closure` |

---

## `closure_proto_compile`

Generates Closure protobuf `.js` files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//closure:deps.bzl", "closure_deps")

closure_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//closure:defs.bzl", "closure_proto_compile")

closure_proto_compile(
    name = "person_closure_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |

---

## `closure_proto_library`

Generates a Closure library with compiled protobuf `.js` files using `closure_js_library` from `rules_closure`

### `WORKSPACE`

```python
load("@build_stack_rules_proto//closure:deps.bzl", "closure_deps")

closure_deps()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//closure:defs.bzl", "closure_proto_library")

closure_proto_library(
    name = "person_closure_library",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### Attributes

| Name | Type | Mandatory | Default | Description |
| ---: | :--- | --------- | ------- | ----------- |
| `deps` | `list<ProtoInfo>` | true | `[]`    | List of labels that provide a `ProtoInfo` (such as `native.proto_library`)          |
| `verbose` | `int` | false | `0`    | The verbosity level. Supported values and results are 1: *show command*, 2: *show command and sandbox after running protoc*, 3: *show command and sandbox before and after running protoc*, 4. *show env, command, expected outputs and sandbox before and after running protoc*          |
