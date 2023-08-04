# cmd/gazelle

This is essentially a copy of the files in bazelbuild/bazel-gazelle/cmd/gazelle.

To upgrade gazelle, one must:

- `go get` the correct version and update `go.mod` file
- `make tidy` to propagate changes to `go.sum` and `go_repositories.bzl`.
- update the version of @bazel_gazelle in `deps/BUILD.bazel`, then run `make
   deps` to regenerate the actual deps tree.
- Compare changes in the source repo to the files here.  It's easiest to just
   copy over each file and see where the diffs are.  Make sure `langs.go`
   includes the `github.com/stackb/rules_proto/language/protobuf`.  Internal
   packages referenced must also be copied over (ugh).  There's probably a more elegant solution to keeping a modified copy of gazelle binary here.
- Since the `proto_gazelle.bzl` rule uses
  `@bazel_gazelle//internal:gazelle.bash.in`, changes there must remain
  compatible with proto_gazelle.  Look at the diff there and make sure the proto_gazelle_impl is satifying the needs of that template.
- Remember that this `cmd/gazelle` must be buildable via the standard go
  toolchain (see proto_repository_tools.bzl):

```py
    args = [
        go_tool,
        "install",
        "-ldflags",
        "-w -s",
        "-gcflags",
        "all=-trimpath=" + env["GOPATH"],
        "-asmflags",
        "all=-trimpath=" + env["GOPATH"],
        "github.com/stackb/rules_proto/cmd/gazelle",
    ]
    result = env_execute(ctx, args, environment = env)
    if result.return_code:
        fail("failed to build tools: " + result.stderr)
```

So any and all deps must be in the vendor tree.
