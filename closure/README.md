# `closure`

| Rule | Description |
| ---: | :--- |
| [closure_proto_compile](#closure_proto_compile) | Generates closure *.js protobuf+gRPC files |
| [closure_proto_library](#closure_proto_library) | Generates a closure_library with compiled protobuf *.js files |

---

## `closure_proto_compile`

Generates closure *.js protobuf+gRPC files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//closure:deps.bzl", "closure_proto_compile")

closure_proto_compile()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//closure:closure_proto_compile.bzl", "closure_proto_compile")

closure_proto_compile(
    name = "person_closure_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile_attrs", "proto_compile_impl")
load("//:aspect.bzl", "ProtoLibraryAspectNodeInfo", "proto_compile_aspect_attrs", "proto_compile_aspect_impl")
load("//:plugin.bzl", "ProtoPluginInfo")

# "Aspects should be top-level values in extension files that define them."

_aspect = aspect(
    implementation = proto_compile_aspect_impl,
    attr_aspects = ["deps"],
    attrs = proto_compile_aspect_attrs + {
        "_plugins": attr.label_list(
            doc = "List of protoc plugins to apply",
            providers = [ProtoPluginInfo],
            default = [
                str(Label("//closure:js")),
            ],
        ),
    },
)

_rule = rule(
    implementation = proto_compile_impl,
    attrs = proto_compile_attrs + {
        "deps": attr.label_list(
            mandatory = True,
            providers = ["proto", ProtoLibraryAspectNodeInfo],
            aspects = [_aspect],
        ),    
    },
)

def closure_proto_compile(**kwargs):
    _rule(
        verbose_string = "%s" % kwargs.get("verbose", 0),
        plugin_options_string = ";".join(kwargs.get("plugin_options", [])),
        **kwargs)

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

## `closure_proto_library`

Generates a closure_library with compiled protobuf *.js files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//closure:deps.bzl", "closure_proto_library")

closure_proto_library()

load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")

closure_repositories(
    omit_com_google_protobuf = True,
)
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//closure:closure_proto_library.bzl", "closure_proto_library")

closure_proto_library(
    name = "person_closure_library",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### `IMPLEMENTATION`

```python
load("//closure:closure_proto_compile.bzl", "closure_proto_compile")
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_js_library")

def closure_proto_library(**kwargs):
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")
    transitive = kwargs.pop("transitive", True)
    transitivity = kwargs.get("transitivity")

    name_pb = name + "_pb"

    closure_proto_compile(
        name = name_pb,
        deps = deps,
        visibility = visibility,
        transitive = transitive,
        transitivity = transitivity,
    )

    closure_js_library(
        name = name,
        srcs = [name_pb],
        deps = ["@io_bazel_rules_closure//closure/protobuf:jspb"],
        visibility = visibility,
        internal_descriptors = [name_pb + "/descriptor.source.bin"],
        lenient = True,
        suppress = [
            "JSC_WRONG_ARGUMENT_COUNT",
        ],
    )
    name = kwargs.get("name")
    deps = kwargs.get("deps")
    visibility = kwargs.get("visibility")

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

