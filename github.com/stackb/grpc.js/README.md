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

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")
closure_repositories(omit_com_google_protobuf = True)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/stackb/grpc.js:closure_grpc_compile.bzl", "closure_grpc_compile")

closure_grpc_compile(
    name = "greeter_grpc.js_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile")

def closure_grpc_compile(**kwargs):
    proto_compile(
        plugins = [
            str(Label("//github.com/stackb/grpc.js:grpc.js")),
        ],
        **kwargs
    )
```

### Mandatory Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| deps   | `list<ProtoInfo>` | `[]`    | List of labels that provide a `ProtoInfo` (`native.proto_library`)          |

### Optional Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| plugins   | `list<ProtoPluginInfo>` | `[]`    | List of labels that provide a `ProtoPluginInfo`          |
| plugin_options   | `list<string>` | `[]`    | List of additional 'global' plugin options (applies to all plugins)          |
| outputs   | `list<generated file>` | `[]`    | List of additional expected generated file outputs          |
| has_services   | `bool` | `False`    | If the proto files(s) have a service rpc, generate grpc outputs          |
| protoc   | `executable file` | `@com_google_protobuf//:protoc`    | The protocol compiler tool          |
| verbose   | `int` | `0`    | 1: *show command*, 2: *show sandbox after*, 3: *show sandbox before*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
| transitive   | `bool` | `False`    | Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies          |

---

## `closure_grpc_library`

Generates protobuf closure library *.js files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/stackb/grpc.js:deps.bzl", "closure_grpc_library")
closure_grpc_library()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")
closure_repositories(omit_com_google_protobuf = True)

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/stackb/grpc.js:closure_grpc_library.bzl", "closure_grpc_library")

closure_grpc_library(
    name = "greeter_grpc.js_library",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
)
```

### `IMPLEMENTATION`

```python
load("//github.com/stackb/grpc.js:closure_grpc_compile.bzl", "closure_grpc_compile")
load("//closure:compile.bzl", "closure_proto_compile")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def closure_grpc_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    verbose = kwargs.get("verbose")
    visibility = kwargs.get("visibility")

    name_pb = name + "_pb"
    name_pb_grpc = name + "_pb_grpc"

    closure_proto_compile(
        name = name_pb,
        deps = deps,
        transitive = True,
        visibility = visibility,
    )
    
    closure_grpc_compile(
        name = name_pb_grpc,
        deps = deps,
        transitive = True,
        visibility = visibility,
        verbose = verbose,
    )

    closure_js_library(
        name = name,
        srcs = [name_pb, name_pb_grpc],
        deps = [
            "@io_bazel_rules_closure//closure/library",
            "@io_bazel_rules_closure//closure/protobuf:jspb",
            "@com_github_stackb_grpc_js//js/grpc/stream:observer",
            "@com_github_stackb_grpc_js//js/grpc/stream/observer:call",
            "@com_github_stackb_grpc_js//js/grpc",
            "@com_github_stackb_grpc_js//js/grpc:api",
            "@com_github_stackb_grpc_js//js/grpc:options",
        ],
        internal_descriptors = [
            name_pb + "/descriptor.source.bin",
            name_pb_grpc + "/descriptor.source.bin",
        ],
        lenient = True,
        suppress = [
            "JSC_WRONG_ARGUMENT_COUNT",
        ],
        visibility = visibility,
    )

```

### Mandatory Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| deps   | `list<ProtoInfo>` | `[]`    | List of labels that provide a `ProtoInfo` (`native.proto_library`)          |

### Optional Attributes

| Name | Type | Default | Description |
| ---: | :--- | ------- | ----------- |
| plugins   | `list<ProtoPluginInfo>` | `[]`    | List of labels that provide a `ProtoPluginInfo`          |
| plugin_options   | `list<string>` | `[]`    | List of additional 'global' plugin options (applies to all plugins)          |
| outputs   | `list<generated file>` | `[]`    | List of additional expected generated file outputs          |
| has_services   | `bool` | `False`    | If the proto files(s) have a service rpc, generate grpc outputs          |
| protoc   | `executable file` | `@com_google_protobuf//:protoc`    | The protocol compiler tool          |
| verbose   | `int` | `0`    | 1: *show command*, 2: *show sandbox after*, 3: *show sandbox before*          |
| include_imports   | `bool` | `True`    | Pass the --include_imports argument to the protoc_plugin          |
| include_source_info   | `bool` | `True`    | Pass the --include_source_info argument to the protoc_plugin          |
| transitive   | `bool` | `False`    | Generated outputs for *.proto directly named in `deps` AND all transitive proto_library dependencies          |

