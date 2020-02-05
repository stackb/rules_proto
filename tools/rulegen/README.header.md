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

Confused about different bazel proto rules?

- [github.com/bazelbuild/rules_proto](https://github.com/bazelbuild/rules_proto):
  The effort to move the native `proto_library` rule to a starlark based one.
  This repo does not define any language-specific rules.
- [github.com/pubref/rules_proto](https://github.com/pubref/rules_proto): The
  predecessor to these rules (1.0).  Re-written to a different organization but
  retains the same primary author @pcj.
- [https://github.com/rules-proto-grpc/rules_proto_grpc](https://https://github.com/rules-proto-grpc/rules_proto_grpc):
  A fork of these rules by @alidell that converts most of the rules to aspects.

