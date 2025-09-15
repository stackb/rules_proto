
## Repository Rules

## proto_repository

From an implementation standpoint, this is very similar to the `go_repository`
rule. Both can download external files and then run the gazelle generator over
the downloaded files. Example:

```python
proto_repository(
    name = "googleapis",
    build_directives = [
        "gazelle:resolve proto google/api/http.proto //google/api:http_proto",
    ],
    build_file_generation = "clean",
    build_file_proto_mode = "file",
    reresolve_known_proto_imports = True,
    proto_language_config_file = "@//:rules_proto_config.yaml",
    strip_prefix = "googleapis-02710fa0ea5312d79d7fb986c9c9823fb41049a9",
    type = "zip",
    urls = ["https://codeload.github.com/googleapis/googleapis/zip/02710fa0ea5312d79d7fb986c9c9823fb41049a9"],
)
```

Takeaways:

- The `urls`, `strip_prefix` and `type` behave similarly to `http_archive`.
- `build_file_proto_mode` is the same the `go_repository` attribute of the same
  name; additionally the value `file` is permitted which generates a
  `proto_library` for each individual proto file.
- `build_file_generation` is the same the `go_repository` attribute of the same
  name; additionally the value `clean` is supported. For example, googleapis
  already has a set of BUILD files; the `clean` mode will remove all the
  existing build files before generating new ones.
- `build_directives` is the same as `go_repository`. Resolve directives in this
  case are needed because the gazelle `language/proto` extension hardcodes a
  proto import like `google/api/http.proto` to resolve to the `@go_googleapis`
  workspace; here we want to make sure that http.proto resolves to the same
  external workspace.
- `proto_language_config_file` is an optional label pointing to a valid
  `config.yaml` file to configure this extension.
- `reresolve_known_proto_imports` is a boolean attribute that has special meaning for
  the googleapis repository. Due to the fact that the builtin gazelle "proto"
  extension has
  [hardcoded logic](https://github.com/bazelbuild/bazel-gazelle/blob/master/language/proto/known_proto_imports.go)
  for what googleapis deps look like, additional work is needed to override
  that. With this sample configuration, the following build command succeeds:

```bash
$ bazel build @googleapis//google/api:annotations_cc_library
Target @googleapis//google/api:annotations_cc_library up-to-date:
  bazel-bin/external/googleapis/google/api/libannotations_cc_library.a
  bazel-bin/external/googleapis/google/api/libannotations_cc_library.so
```

Another example using the Bazel repository:

```python
load("@build_stack_rules_proto//rules/proto:proto_repository.bzl", "proto_repository")

proto_repository(
    name = "bazelapis",
    build_directives = [
        "gazelle:exclude third_party",
        "gazelle:proto_language go enable true",
        "gazelle:proto_language closure enabled true",
        "gazelle:prefix github.com/bazelbuild/bazelapis",
    ],
    build_file_expunge = True,
    build_file_proto_mode = "file",
    cfgs = ["//proto:config.yaml"],
    imports = [
        "@googleapis//:imports.csv",
        "@protobufapis//:imports.csv",
        "@remoteapis//:imports.csv",
    ],
    strip_prefix = "bazel-02ad3e3bc6970db11fe80f966da5707a6c389fdd",
    type = "zip",
    urls = ["https://codeload.github.com/bazelbuild/bazel/zip/02ad3e3bc6970db11fe80f966da5707a6c389fdd"],
)
```

Takeaways:

- The `build_directives` are setting the `gazelle:prefix` for the `language/go`
  plugin; two `proto_language` configs named in the `proto/config.yaml` are
  being enabled.
- `build_file_expunge` means _remove all existing BUILD files before generating
  new ones_.
- `build_file_proto_mode = "file"` means _make a separate `proto_library` rule
  for every `.proto` file_.
- `cfgs = ["//proto:config.yaml"]` means _use the configuration in this YAML
  file as a base set of gazelle directives_. It is easier/more consistent to
  share the same `config.yaml` file across multiple `proto_repository` rules.

The last one `imports = ["@googleapis//:imports.csv", ...]` requires extra
explanation. When the `proto_repository` gazelle process finishes, it writes a
file named `imports.csv` in the root of the external workspace. This file
records the association between import statements and bazel labels, much like a
`gazelle:resolve` directive:

```csv
# GENERATED FILE, DO NOT EDIT (created by gazelle)
# lang,imp.lang,imp,label
go,go,github.com/googleapis/gapic-showcase/server/genproto,@googleapis//google/example/showcase/v1:compliance_go_proto
go,go,google.golang.org/genproto/firestore/bundle,@googleapis//google/firestore/bundle:bundle_go_proto
go,go,google.golang.org/genproto/googleapis/actions/sdk/v2,@googleapis//google/actions/sdk/v2:account_linking_go_proto
```

Therefore, the `imports` attribute assists gazelle in figuring how to resolve
imports. Therefore, when gazelle is preparing a `go_library` rule and finds a
`main.go` file having an import on
`google.golang.org/genproto/firestore/bundle`, it knows to put
`@googleapis//google/firestore/bundle:bundle_go_proto` in the rule `deps`.

To take advantage of this mechanism in the default workspace, use the
`proto_gazelle` rule.
