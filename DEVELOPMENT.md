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
- `//cmd/depsgen`: a helper tool used by the `depsgen` rule (see `//deps:*`)
  that generates `.bzl` files based on the transitive set of `proto_repository`
  rules.
- `//cmd/examplegen`: a helper tool that generates markdown and a `_test.go`
  file used by the `gazelle_testdata_example` rule (see `//example/golden:*`).
- `//cmd/gencopy`: a helper tool that is used to copy generated files back into
  the source tree.
- `//example/...`: example gazelle directives and rules.
- `//deps:*`: contains the `//deps:BUILD.bazel` files where all project
  dependencies are declared as `proto_repository` rules.  These are linked
  together via their `deps` attribute.  Each `depsgen(name = "foo")` rule in
  that file is used to create the corresponding `deps/foo_deps.bzl`, which is
  generated and then copied back into the source tree.  Use `make deps` target
  to regenerate these files.
- `go.mod, go.sum`.  The are the inputs to generate the `go_deps.bzl` file.  Use
  `bazel run //:update_go_deps` to regenerate this file.
- `//language/protobuf`: contains the `protobuf` language extension.  The label
  `@build_stack_rules_proto//language/protobuf` should be referenced in the
  `gazelle_binary` rule.
- `//language/example` contains a no-op gazelle extension that can be used an
  empty template.
- `//pkg/...`: contains all go packages for tools.
- `//pkg/language/protobuf`: contains the actual protobuf gazelle extension
  implementation.
- `//pkg/language/noop`: contains an empty extension that can be used as a starting template.
- `//pkg/protoc`: a core package that defines the registry, language, rule
  interfaces, and much of the implementation.  I'm not sure why I named it
  `protoc`.
- `//pkg/rule/...`: rule implementations.
- `//plugin/...`: directories here contain `proto_plugin` rules.  In some cases,
  the plugin tool is built here, possibly depending on language-specific deps.
  For example, the `//plugin/stephenh/ts-proto` tool requires NPM dependencies.
- `//third_party/...`: third party BUILD rules.  TODO: get rid of this.
- `//toolchain`.  The toolchain definition used by the `proto_compile` rule.

## Vendoring Files

The `vendor/` tree is used by the standard go build toolchain when constructing
the binary for the `proto_repository` rule.  If new dependencies are added to
the tool, they may need to be also added to `vendor/`.  Use `make go_mod_vendor`
for this.  That target will run the vendor command and then remove any
`BUILD.bazel` files.  This is needed because the `//:all_files` target attempts
to glob up everything in the `vendor/` directory.  This does not work correctly
if there are `BUILD.bazel` files in vendor.  These files are also listed in the
`.gitignore` file.