# py_proto_compile

## `WORKSPACE`

```python
load("@build_stack_rules_proto//python:py_proto_compile.bzl", "py_proto_compile")

py_proto_compile()
```

## `BUILD.bazel`

```python
load("@build_stack_rules_proto//python:py_proto_compile.bzl", "py_proto_compile")

py_proto_compile(
    name = "py_proto_compile_person_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

## Plugins

| Label | Tool | Outputs |
| ---- | ---- | ------- |
| `//python:python_plugin` |  |  `{protopath|python}_pb2.py` |


## Dependencies

load("@build_stack_rules_proto//python:py_proto_compile.bzl", "py_proto_compile")

py_proto_compile()
