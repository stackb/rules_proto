# Development Notes

## Project Layout

This repo attempts to keep language-specific dependencies separate, meaning if
you use `stackb/rules_proto` for C++, you shouldn't need to transitively pull in
scala deps, for example.

- `//rules:*all*`: core starlark rules and providers.  Bzl files here should not
  load any language-specific things.
- `//rules/{LANG}:*all*`: language-specific rules.  Bzl files here may reference
  language-specific deps.
- `//cmd/gazelle`: a copy of the gazelle frontend with some minor modifications.
  The `langs.go` file is added which loads the protobuf language.  This copy of
  gazelle ends up being built by the standard go toolchain for the
  `proto_repository` rule.
- `//cmd/examplegen`: a helper tool that generates markdown and a `_test.go`
  file used by the `gazelle_testdata_example` rule (see `//example/golden:*`).
- `//cmd/gencopy`: a helper tool that is used to copy generated files back into
  the source tree.
- `//example/...`: example gazelle directives and rules.
- `//language/protobuf`: contains the `protobuf` language extension.  The label
  `@build_stack_rules_proto//language/protobuf` should be referenced in the
  `gazelle_binary` rule.
- `//language/example` contains a no-op gazelle extension that can be used an
  empty template.
- `//pkg/...`: contains all go packages for tools.
- `//pkg/language/protobuf`: contains the actual protobuf gazelle extension
  implementation.
- `//pkg/language/noop`: contains an empty extension that can be used as a
  starting template.
- `//pkg/protoc`: a core package that defines the registry, language, rule
  interfaces, and much of the implementation.
- `//pkg/rule/...`: rule implementations.
- `//plugin/...`: directories here contain `proto_plugin` rules.  In some cases,
  the plugin tool is built here, possibly depending on language-specific deps.
  For example, the `//plugin/stephenh/ts-proto` tool requires NPM dependencies.
- `//third_party/...`: third party BUILD rules.  TODO: get rid of this.
- `//toolchain`.  The toolchain definition used by the `proto_compile` rule.

## Vendoring Files

The `vendor/` tree is used by the standard go build toolchain when constructing
the binary for the `proto_repository` rule.  If new dependencies are added to
the tool, they may need to be also added to `vendor/`.  Use `make tidy` for
this.  That target will run the vendor command and then remove any `BUILD.bazel`
files.  This is needed because the `//:all_files` target attempts to glob up
everything in the `vendor/` directory.  This does not work correctly if there
are `BUILD.bazel` files in vendor.  These files are also listed in the
`.gitignore` file.

## Debugging `gazelle_testdata_example` tests.

The target `//example/golden:golden_test` runs a test foreach subdirectory of
`example/golden/testdata/`.  These tests stage all the `BUILD.in` files,
optionally using a `.gazelle.args` file to run gazelle.  It then compares the
actual BUILD files produced to the `BUILD.out` files.  These are generally
straightforward to debug as a diff is printed out if the test fails.

Use `bazel query //example/golden:all` to see all the actual tests in this
package.

The tests defined by `gazelle_testdata_example` glob up one of the
`testdata/scala/` directories and generates a `_test.go` file.  Use `bazel build
//example/golden:scala` to print out the generated test file (`cat
bazel-bin/example/golden/scala_test.go`).

This type of test packages up the `BUILD.out` files and uses `go_bazel_test` to
actually run `bazel build ...` on it.  So, the `:golden_test` target is used to
assert that the correct gazelle output is generated, while `:scala_test` is used
to assert that the outbut actually builds.

The effective `MODULE.bazel` file is a concatenation of the `workspace_template`
(e.g. `default.MODULE.bazel`) and the one in the testdata example dir.

One way to debug these kinds of tests is to open an editor in the effective
directory where the test is constructed (this will be printed out if it fails).
For example `(cd
/private/var/tmp/_bazel_i868039/7b08591af2b3d71f45f2e4029050db37/bazel_testing/bazel_go_test/main
&& code .)`.  You can then use `bazel build` directly in this space to explore
what's going on.

## Repository input tracking

`proto_repository` watches its `cfgs` YAML files and `imports` CSV files.
Gazelle also writes a temporary JSON manifest through
`-proto_config_inputs_out` listing the Starlark plugin and rule files loaded
from YAML or `-proto_plugin` / `-proto_rule`. The repository rule watches those
files and removes the manifest. Sources inside the fetched repository are
already covered by the repository definition and are not watched separately.
The repository Gazelle tool also watches its Go sources so source changes
rebuild the executable.

The following regression test uses real Bazel and Gazelle in a temporary
workspace. It checks configuration, CSV, and Starlark edits, propagation to a
consumer repository, and reuse when no inputs change. Archives and dependencies
are supplied locally; the test does not require network access.

```sh
go build -mod=vendor -o /tmp/rules-proto-gazelle ./cmd/gazelle
python3 tools/test_proto_repository_inputs.py \
  --bazel "$(command -v bazel)" \
  --gazelle /tmp/rules-proto-gazelle \
  --gazelle-repo "$(bazel info output_base)/external/gazelle+"
```

`--gazelle-repo` must point to the fetched bazel-gazelle sources; adjust the
canonical repository name if your Bazel version uses a different spelling.
