# Writing Custom Plugins and Rules

Custom plugin implementations and rule implementations can be written in golang
or starlark.  Golang implementations are statically compiled into the final
`gazelle_binary` whereas starlark plugins are evaluated at gazelle runtime.

## +/- of golang implementations

- `+` Full power of a statically compiled language, the golang stdlib, and
  external dependencies.
- `+` Easier to test.
- `+` API not experimental.
- `-` Cannot be used in a `proto_repository` rule without forking
  stackb/rules_proto.
- `-` Initial setup harder, often housed within your own custom gazelle
  extension.

Until a dedicated tutorial is available, please consult the source code for
examples.

## +/- of starlark implementations

- `+` More familiar to developer with starlark experience but not golang.
- `+` Easier setup (*.star files in your gazelle repository)
- `+` Possible to use in conjunction with the `proto_repository` rule.
- `-` Limited API: can only reference state that has been already configured via gazelle directives.
- `-` Not possible to implement stateful design.
- `-` No standard library.

Until a dedicated tutorial is available, please consult the reference example in
`example/testdata/starlark_java`.
