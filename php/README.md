# `php`

| Rule | Description |
| ---: | :--- |
| [php_proto_compile](#php_proto_compile) | Generates php protobuf artifacts |
| [php_grpc_compile](#php_grpc_compile) | Generates php protobuf+gRPC artifacts |

---

## `php_proto_compile`

Generates php protobuf artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//php:deps.bzl", "php_proto_compile")

php_proto_compile()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//php:php_proto_compile.bzl", "php_proto_compile")

php_proto_compile(
    name = "person_php_proto",
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
                str(Label("//php:php")),
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

def php_proto_compile(**kwargs):
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

## `php_grpc_compile`

Generates php protobuf+gRPC artifacts

### `WORKSPACE`

```python
load("@build_stack_rules_proto//php:deps.bzl", "php_grpc_compile")

php_grpc_compile()

load("@com_github_grpc_grpc//bazel:grpc_deps.bzl", "grpc_deps")

grpc_deps()
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//php:php_grpc_compile.bzl", "php_grpc_compile")

php_grpc_compile(
    name = "greeter_php_grpc",
    deps = ["@build_stack_rules_proto//example/proto:greeter_grpc"],
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
                str(Label("//php:php")),
                str(Label("//php:grpc_php")),
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

def php_grpc_compile(**kwargs):
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

