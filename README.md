# `rules_proto`

Bazel skylark rules for building protocol buffers +/- gRPC :sparkles:.

<table border="0"><tr>
<td><img src="https://bazel.build/images/bazel-icon.svg" height="180"/></td>
<td><img src="https://github.com/pubref/rules_protobuf/blob/master/images/wtfcat.png" height="180"/></td>
<td><img src="https://avatars2.githubusercontent.com/u/7802525?v=4&s=400" height="180"/></td>
</tr><tr>
<td>Bazel</td>
<td>rules_proto</td>
<td>gRPC</td>
</tr></table>

Confused about the different bazel proto rules?

- [github.com/pubref/rules_proto](https://github.com/pubref/rules_proto): The
  predecessor to stackb/rules_proto, now archived and not under development.
- [github.com/stackb/rules_proto](https://github.com/stackb/rules_proto): The
  pubref ruleset was re-written basically from scratch, using a set of golang
  based templates to generate the actual rules, with generic `plugin.bzl`,
  `compile.bzl`, and `aspect.bzl` abstractions to support most protoc plugins.
  Supporting a great deal of languages was a useful means to hone the
  abstractions, but with so many interlinked dependencies, maintainance was
  troublesome, and led to maintainer burnout.
- [github.com/rules-proto-grpc/rules_proto_grpc](https://github.com/rules-proto-grpc/rules_proto_grpc):
  the successor to stackb/rules_proto. @alidell drove the initial `aspect.bzl`
  implementation from its original prototype, got it working with all the rules,
  and forked the repository.  @aliddell is doing a good job maintaining that
  repository.
- [github.com/bazelbuild/rules_proto](https://github.com/bazelbuild/rules_proto):
  The native `proto_library` is of course implemented in java within the
  `bazelbuild/bazel` repository.  That rules_proto repo is primarily concerned
  with transitioning the native rule to a starlark-based one in it's own
  repository.

So what is `stackb/rules_proto` now?  These are the goals:

* Goal 1: Rather that providing the actual rule implementations, **simplify the
  repository and provide a set of working examples**.  In this manner
  `stackb/rules_proto` aims to provide a "thin veneer" over existing proto rules
  across the bazel ecosystem and "smooth over" the rough corners.  In this way
  it acts as a sort of "bazel federation" for proto rules; a place where a users
  can refer to a set of dependencies that are known to work together.

* Goal 2: **not duplicate existing proto rules**.  For example, although it is
  academically interesting to provide an independent `go_proto_library` rule,
  @jayconrad has done a solid job implementing this rule already and has a good
  handle on the corner cases in the go toolchain.  Therefore, while
  `rules_proto` provides `go_proto_library` and `go_grpc_library` rules, it is
  simply exporting the original `rules_go` implementation, and showing how to
  setup the dependencies.  Same goes for `cc_grpc_library` and
  `java_grpc_library`.

* Goal 3: **not duplicate existing depdendency loading schemes**.  Rather than
  setup `@zlib`, `@boringssl`, the myriad `@bind`s here, lean on
  `bazelbuild/rules_proto` for esablishing a `com_google_protobuf`,
  `com_github_grpc_grpc` for `grpc_deps`, `io_bazel_grpc_grpc_java`, etc...

As always, contributions are always welcome!

- Feb 2020