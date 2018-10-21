# `rules_proto` [![Build Status](https://travis-ci.org/pubref/rules_proto.svg?branch=master)](https://travis-ci.org/stackb/rules_proto)

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

These rules are the successor to <https://github.com/pubref/rules_protobuf> and
are in a pre-release status.  The primary goals are:

1. Interoperate with the native `proto_library` rules and other proto support in
   the bazel ecosystem as much as possible.
2. Provide a `proto_plugin` rule to support custom protoc plugins.
3. Minimal dependency loading.  Proto rules should not pull in more dependencies
   than they absolutely need.

> NOTE: These rules are in a *pre-release* state.  The routeguide examples have
> been developed thus far with the goal of getting them to compile/build, not to
> ensure the routeguide client/server is actually correct.  Do expect everything
> to compile, but not to work right!

## Usage

Add rules_proto your `WORKSPACE`:

```python
load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
    name = "build_stack_rules_proto",
    urls = ["https://github.com/stackb/rules_proto/archive/{{ .Ref }}.tar.gz"],
    sha256 = {{ .Sha256 }},
    strip_prefix = "rules_proto-{{ .Ref }}",
)
```

Follow instructions in the language-specific `README.md` for additional
workspace dependencies and build rule usage.  
