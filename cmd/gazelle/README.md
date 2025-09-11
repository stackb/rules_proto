# cmd/gazelle

This is essentially a copy of the files in bazel-contrib/bazel-gazelle/cmd/gazelle.

To upgrade gazelle, one must:

- Compare changes in the source repo@version to the files here.  It's easiest to
  just copy over each file and see where the diffs are.  Make sure `langs.go`
  includes `github.com/stackb/rules_proto/language/protobuf`.
- Since the `proto_gazelle.bzl` rule uses
  `@bazel_gazelle//internal:gazelle.bash.in`, changes there must remain
  compatible with proto_gazelle.  Look at the diff there and make sure the
  `_proto_gazelle_impl` is satifying the needs of that template.
- Remember that this `cmd/gazelle` must be buildable via the standard go
  toolchain (see proto_repository_tools.bzl). All deps must be in the vendor
  tree.
