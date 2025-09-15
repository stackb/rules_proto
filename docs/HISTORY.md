
# History

## v0

The original rules_proto was <https://github.com/pubref/rules_proto>. This was
redesigned around the `proto_library` rule and moved to
<https://github.com/stackb/rules_proto>.  

## v1

These were pretty good, but open source maintainance was an issue trying to keep
up with the matrix of (language * dependencies).

Partially as a result of that, this ruleset was forked to
<https://github.com/rules-proto-grpc/rules_proto_grpc>. Aspect-based compilation
was featured for quite a while there but later removed.  Adam does a great job
maintaining that repo.  If you're not using gazelle, check it out!

## v2

`stackb/rules_proto` adopted gazelle as a primary means to manage proto rules.
This shifted the burden of dependency management primarily to the consuming
repos where it belongs.

The gazelle based design makes some things simpler because the **content of the
proto files is the source of truth**.  Due to the fact that Bazel does not
permit reading/interpreting a file during the scope of an action, it is
impossible to make a decision about what to do. A prime example of this is the
`go_package` option. If the `go_package` option is present, the location of the
output file for `protoc-gen-go` is completely different. As a result, the
information about the go_package metadata ultimately needs to be duplicated so
that the build system can know about it.

The gazelle-based approach moves all that messy interpretation and evaluation
into a precompiled state; as a result, the work that needs to be done in the
action itself is simplified.

Furthermore, by parsing the proto files it is easy to support complex custom
plugins that do things like:

- Emit no files (only assert/lint).
- Emit a file only if a specific enum constant is found. With the previous
  design, this was near impossible. With the `v2` design, the `protoc.Plugin`
  implementation can trivially perform that evaluation because it is handed the
  complete proto AST during gazelle evaluation.

## v3

The way gazelle stored resolve data in memory changed in 0.35.  Uprading to this
in https://github.com/stackb/rules_proto/pull/357 was a breaking change.

## v4

With v4, `WORKSPACE` support was dropped in favor of `bzlmod`.
