# rules/go

## `proto_go_modules`

The `proto_go_modules` rule is used to vendor generated golang protobuf assets
from external workspaces into the default workspace.

For example, imagine you are building an application that requires the remote
execution API.  The top-level protobuf package
`build/bazel/remote/execution/v2/remote_execution.proto` has a fairly complex
proto dependency tree including protos from the googleapis and protobuf
well-known protos.

One solution is to manually copy all the needed `.proto` files from the
different repos into the default workspace and generate protos from there.  This
is a common solution but can be troublesome to maintain.

The solution described here (`bazel run @//:proto_go_modules`) has the following
two side effects: 

1. creates a "vendored" file tree in `./local`
2. modifies the `go.mod` file with new `replace` directives

```sh
$ tree ./local
local
├── github.com
│   └── bazelbuild
│       └── remoteapis
│           └── build
│               └── bazel
│                   ├── remote
│                   │   └── execution
│                   │       └── v2
│                   │           └── remote_execution
│                   │               ├── go.mod
│                   │               ├── remote_execution_grpc.pb.go
│                   │               └── remote_execution.pb.go
│                   └── semver
│                       └── semver
│                           ├── go.mod
│                           └── semver.pb.go
└── google.golang.org
    ├── genproto
    │   └── googleapis
    │       ├── api
    │       │   └── annotations
    │       │       ├── annotations.pb.go
    │       │       ├── client.pb.go
    │       │       ├── field_behavior.pb.go
    │       │       ├── go.mod
    │       │       ├── http.pb.go
    │       │       └── resource.pb.go
    │       ├── longrunning
    │       │   ├── go.mod
    │       │   ├── operations_grpc.pb.go
    │       │   └── operations.pb.go
    │       └── rpc
    │           └── status
    │               ├── go.mod
    │               └── status.pb.go
    └── protobuf
        └── types
            ├── descriptorpb
            │   ├── descriptor.pb.go
            │   └── go.mod
            └── known
                ├── anypb
                │   ├── any.pb.go
                │   └── go.mod
                ├── durationpb
                │   ├── duration.pb.go
                │   └── go.mod
                ├── emptypb
                │   ├── empty.pb.go
                │   └── go.mod
                ├── timestamppb
                │   ├── go.mod
                │   └── timestamp.pb.go
                └── wrapperspb
                    ├── go.mod
                    └── wrappers.pb.go
```

```bash
module github.com/stackb/bazel-aquery-differ

go 1.23.1

require (
    ...
)

replace (
    github.com/bazelbuild/remoteapis/build/bazel/semver/semver => ./local/github.com/bazelbuild/remoteapis/build/bazel/semver/semver
    google.golang.org/protobuf/types/descriptorpb => ./local/google.golang.org/protobuf/types/descriptorpb
    google.golang.org/genproto/googleapis/api/annotations => ./local/google.golang.org/genproto/googleapis/api/annotations
    google.golang.org/protobuf/types/known/anypb => ./local/google.golang.org/protobuf/types/known/anypb
    google.golang.org/genproto/googleapis/rpc/status => ./local/google.golang.org/genproto/googleapis/rpc/status
    google.golang.org/protobuf/types/known/durationpb => ./local/google.golang.org/protobuf/types/known/durationpb
    google.golang.org/protobuf/types/known/timestamppb => ./local/google.golang.org/protobuf/types/known/timestamppb
    google.golang.org/protobuf/types/known/emptypb => ./local/google.golang.org/protobuf/types/known/emptypb
    google.golang.org/genproto/googleapis/longrunning => ./local/google.golang.org/genproto/googleapis/longrunning
    google.golang.org/protobuf/types/known/wrapperspb => ./local/google.golang.org/protobuf/types/known/wrapperspb
    github.com/bazelbuild/remoteapis/build/bazel/remote/execution/v2/remote_execution => ./local/github.com/bazelbuild/remoteapis/build/bazel/remote/execution/v2/remote_execution
)
```

The `local/` dir should be checked into `git`.  This has the following benefits:

- `.pb.go` files are trivially available.  This keeps the IDE happy, and tool
  operations like `go mod tidy` work as expected.
- The source of truth remains explicitly specified and can be easily updated as
  needed.

The remainder of this document goes through this step-by-step.  For the concrete
example, see the repo <https://github.com/stackb/bazel-aquery-differ/>.

### Step 1: configure `proto_repository` rules in `MODULE.bazel`

```py
module(
    name = "bazel-aquery-differ",
    version = "0.0.0",
)

# -------------------------------------------------------------------
# Bazel Dependencies
# -------------------------------------------------------------------

bazel_dep(name = "rules_proto", version = "7.1.0")
bazel_dep(name = "rules_go", version = "0.57.0")
bazel_dep(name = "gazelle", version = "0.45.0")
bazel_dep(name = "build_stack_rules_proto", version = "0.0.0")

# -------------------------------------------------------------------
# Configuration: Protobuf Deps
# -------------------------------------------------------------------

proto_repository = use_extension("@build_stack_rules_proto//extensions:proto_repository.bzl", "proto_repository", dev_dependency = True)
proto_repository.archive(
    name = "googleapis",
    build_directives = [
        "gazelle:proto_go_modules_target_dir ROOT",
        "gazelle:exclude google/example",
        "gazelle:exclude google/ads/googleads/v7/services",
        "gazelle:exclude google/ads/googleads/v8/services",
        "gazelle:proto_language go enabled true",
    ],
    build_file_generation = "clean",
    build_file_proto_mode = "file",
    cfgs = ["//:rules_proto_config.yaml"],
    imports = ["@protoapis//:imports.csv"],
    reresolve_known_proto_imports = True,
    sha256 = "95da12951c7d570980d5152f6cca9e1cb795ddc6b6dd7e9423bdffde28290f7a",
    strip_prefix = "googleapis-02710fa0ea5312d79d7fb986c9c9823fb41049a9",
    type = "zip",
    urls = [
        "https://codeload.github.com/googleapis/googleapis/zip/02710fa0ea5312d79d7fb986c9c9823fb41049a9",
    ],
)
proto_repository.archive(
    name = "remoteapis",
    build_directives = [
        "gazelle:proto_go_modules_target_dir ROOT",
        "gazelle:exclude third_party",
        "gazelle:exclude build/bazel/remote/asset/v1",
        "gazelle:exclude build/bazel/remote/logstream/v1",
        "gazelle:proto_language go enable true",
        "gazelle:proto_plugin protoc-gen-go option Mbuild/bazel/remote/execution/v2/remote_execution.proto=github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2",
        "gazelle:proto_plugin protoc-gen-go option Mbuild/bazel/semver/semver.proto=github.com/bazelbuild/remote-apis/build/bazel/semver",
        "gazelle:proto_plugin protoc-gen-go-grpc option Mbuild/bazel/remote/execution/v2/remote_execution.proto=github.com/bazelbuild/remote-apis/build/bazel/remote/execution/v2",
        "gazelle:proto_plugin protoc-gen-go-grpc option Mbuild/bazel/semver/semver.proto=github.com/bazelbuild/remote-apis/build/bazel/semver",
    ],
    build_file_generation = "clean",
    build_file_proto_mode = "file",
    cfgs = ["//:rules_proto_config.yaml"],
    imports = [
        "@googleapis//:imports.csv",
        "@protoapis//:imports.csv",
    ],
    sha256 = "743d2d5b5504029f3f825beb869ce0ec2330b647b3ee465a4f39ca82df83f8bf",
    strip_prefix = "remote-apis-636121a32fa7b9114311374e4786597d8e7a69f3",
    type = "zip",
    urls = [
        "https://codeload.github.com/bazelbuild/remote-apis/zip/636121a32fa7b9114311374e4786597d8e7a69f3",
    ],
)
proto_repository.archive(
    name = "protoapis",
    build_directives = [
        "gazelle:proto_go_modules_target_dir ROOT",
        "gazelle:exclude testdata",
        "gazelle:exclude google/protobuf/compiler/ruby",
        "gazelle:exclude google/protobuf/util/internal/testdata",
        "gazelle:proto_language go enable true",
    ],
    build_file_generation = "clean",
    build_file_proto_mode = "file",
    cfgs = ["//:rules_proto_config.yaml"],
    deleted_files = [
        "google/protobuf/*test*.proto",
        "google/protobuf/*unittest*.proto",
        "google/protobuf/compiler/cpp/*test*.proto",
        "google/protobuf/util/*test*.proto",
        "google/protobuf/util/*unittest*.proto",
        "google/protobuf/util/json_format*.proto",
    ],
    sha256 = "4514213c25a5b87e1948aeeb4c40effc55d11d60871ca5b903a2779005fc48ce",
    strip_prefix = "protobuf-9650e9fe8f737efcad485c2a8e6e696186ae3862/src",
    type = "zip",
    urls = [
        "https://codeload.github.com/protocolbuffers/protobuf/zip/9650e9fe8f737efcad485c2a8e6e696186ae3862",
    ],
)
use_repo(
    proto_repository,
    "googleapis",
    "protoapis",
    "remoteapis",
)
```

Each one has a base YAML configuration given in `rules_proto_config.yaml` and a
few extra gazelle directives that are injected into the generated root
`BUILD.bazel` file that:

1. turns on go protobuf rule generation for `go` (`gazelle:proto_language go
   enable true`)
2. turns on the gazelle lang `proto_go_modules` by setting the target package to
   the special token `ROOT` (`gazelle:proto_go_modules_target_dir ROOT`).

When the requested (e.g. `bazel query @protoapis//...`), bazel will:

1. download and extract the zip file
2. remove all existing `BUILD.bazel` files in the workspace
3. generate a new `BUILD.bazel` file in the external root containing the
   `build_directives`
4. remove any additional files / glob patterns named in `deleted_files`
5. run `gazelle` in the external workspace, generating new `BUILD.bazel` files.
   (note: the `gazelle` binary used is not the "normal" gazelle from
   `@bazel_gazelle` but the one from `@build_stack_rules_proto`, which has
   additional extensions pre-packaged).

The generated `proto_go_modules` rule appears as follows.  The `deps` of this
rule include all the generated `proto_go_library` rules in the repo (`GoArchive`
providers):

```sh
bazel query @protoapis//:proto_go_modules --output build
```

```py
proto_go_modules(
    name = "proto_go_modules",
    visibility = ["//visibility:public"],
    deps = [
        "@protoapis//google/protobuf:any_go_proto",
        "@protoapis//google/protobuf:api_go_proto",
        "@protoapis//google/protobuf:descriptor_go_proto",
        "@protoapis//google/protobuf:duration_go_proto",
        "@protoapis//google/protobuf:empty_go_proto",
        "@protoapis//google/protobuf:field_mask_go_proto",
        "@protoapis//google/protobuf:source_context_go_proto",
        "@protoapis//google/protobuf:struct_go_proto",
        "@protoapis//google/protobuf:timestamp_go_proto",
        "@protoapis//google/protobuf:type_go_proto",
        "@protoapis//google/protobuf:wrappers_go_proto",
        "@protoapis//google/protobuf/compiler:plugin_go_proto",
    ],
)
```

### Step 2: manually author a `proto_go_modules` rule in the default workspace

```py
load("@build_stack_rules_proto//rules/go:proto_go_modules.bzl", "proto_go_modules")

proto_go_modules(
    name = "proto_go_modules",
    go_version = "go.mod",
    imports = [
        "github.com/bazelbuild/remoteapis/build/bazel/remote/execution/v2/remote_execution",
    ],
    modules = [
        "@googleapis//:proto_go_modules",
        "@protoapis//:proto_go_modules",
        "@remoteapis//:proto_go_modules",
    ],
)
```

This rule is similar to the one(s) generated in external workspace (via the
`proto_go_modules` gazelle extension), but has a few differences:

- the `modules` attribute names other `proto_go_modules` in the external
  `proto_repository` workspaces (having provider `ProtoGoModulesInfo`).  This
  sets up a "universe" of go_library rules that are uniquely identified by their
  `importpath`.
- the `imports` attribute names a list of desired "top-level" importpaths that
  should be copied/vendored.  The rule will transitively identify all child
  proto-related packages and copy their source files into the default workspace.
  1. `./local/{importpath}/{basename}.pb.go` is the destination file pattern for sources.
  2. `./local/{importpath}/go.mod` is generated for the module.
- the `go_version` attribute specifies the go version for the generated `go.mod`
  files.  It takes one of two forms:
  1. case `go_version = "go_mod"`: if the rule attribute has the special value
  `go.mod`, the go version will be taken from the root `go.mod` file in the
  default workspace.
  1. otherwise, the version specified will be used (e.g. `go_version =
     "1.23.0"`).

### Step 3: execute the shell script generated by the `proto_go_modules` rule

To generate the `local/` dir and update the `go.mod` file:

```sh
$ bazel run //:proto_go_modules
```

If you inspect the script generated by `bazel build //:proto_go_modules`, it
looks like:

```sh
% cat bazel-bin/proto_go_modules
set -euox pipefail
cwd=$PWD
cd $BUILD_WORKING_DIRECTORY
go_version=$(grep '^go' < go.mod)

# module=github.com/bazelbuild/remoteapis/build/bazel/semver/semver
mkdir -p ./local/github.com/bazelbuild/remoteapis/build/bazel/semver/semver
echo 'module github.com/bazelbuild/remoteapis/build/bazel/semver/semver' > ./local/github.com/bazelbuild/remoteapis/build/bazel/semver/semver/go.mod
echo "${go_version}" >> ./local/github.com/bazelbuild/remoteapis/build/bazel/semver/semver/go.mod
cp -f "${cwd}/../build_stack_rules_proto++proto_repository+remoteapis/build/bazel/semver/semver.pb.go" ./local/github.com/bazelbuild/remoteapis/build/bazel/semver/semver/semver.pb.go
echo '# ../build_stack_rules_proto++proto_repository+remoteapis/build/bazel/semver/semver.pb.go' >> ./local/github.com/bazelbuild/remoteapis/build/bazel/semver/semver/go.mod
go mod edit -replace github.com/bazelbuild/remoteapis/build/bazel/semver/semver=./local/github.com/bazelbuild/remoteapis/build/bazel/semver/semver
```

> NOTE: the script does not remove old entries, so you may need to `rm -rf
> ./local` and manually tweak the `go.mod` file in the case of removing
> `imports` from the `proto_go_modules` rule.