# `ts-protoc-gen`

| Rule | Description |
| ---: | :--- |
| [ts_proto_compile](#ts_proto_compile) | Generates typescript protobuf t.ds files |

---

## `ts_proto_compile`

> NOTE: this rule is EXPERIMENTAL.  It may not work correctly or even compile!

Generates typescript protobuf t.ds files

### `WORKSPACE`

```python
load("@build_stack_rules_proto//github.com/improbable-eng/ts-protoc-gen:deps.bzl", "ts_proto_compile")
ts_proto_compile()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()

load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

load("@org_pubref_rules_node//node:rules.bzl", "node_repositories", "yarn_modules")
node_repositories()

load("@build_bazel_rules_nodejs//:defs.bzl", "node_repositories")
node_repositories(package_json = ["@ts_protoc_gen//:package.json"])

load("@build_bazel_rules_typescript//:defs.bzl", "ts_setup_workspace")
ts_setup_workspace()

load("@io_bazel_rules_webtesting//web:repositories.bzl", "browser_repositories", "web_test_repositories")
web_test_repositories()

load("@build_bazel_rules_nodejs//:defs.bzl", "npm_install")
npm_install(
    name = "deps",
    package_json = "@ts_protoc_gen//:package.json",
    package_lock_json = "@ts_protoc_gen//:package-lock.json",
)
```

### `BUILD.bazel`

```python
load("@build_stack_rules_proto//github.com/improbable-eng/ts-protoc-gen:ts_proto_compile.bzl", "ts_proto_compile")

ts_proto_compile(
    name = "person_ts-protoc-gen_proto",
    deps = ["@build_stack_rules_proto//example/proto:person_proto"],
)
```

### `IMPLEMENTATION`

```python
load("//:compile.bzl", "proto_compile_attrs", "proto_compile_impl")
load("//:aspect.bzl", "ProtoLibraryAspectNodeInfo", "proto_compile_aspect_attrs", "proto_compile_aspect_impl")
load("//:plugin.bzl", "ProtoPluginInfo")

# "Aspects should be top-level values in extension files that define them."

ts_proto_compile_aspect = aspect(
    implementation = proto_compile_aspect_impl,
    provides = ["proto_compile", ProtoLibraryAspectNodeInfo],
    attr_aspects = ["deps"],
    attrs = proto_compile_aspect_attrs + {
        "_plugins": attr.label_list(
            doc = "List of protoc plugins to apply",
            providers = [ProtoPluginInfo],
            default = [
                str(Label("//github.com/improbable-eng/ts-protoc-gen:ts")),
            ],
        ),
    },
)

_rule = rule(
    implementation = proto_compile_impl,
    attrs = proto_compile_attrs + {
        "deps": attr.label_list(
            mandatory = True,
            providers = ["proto", "proto_compile", ProtoLibraryAspectNodeInfo],
            aspects = [ts_proto_compile_aspect],
        ),    
    },
)

def ts_proto_compile(**kwargs):
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

