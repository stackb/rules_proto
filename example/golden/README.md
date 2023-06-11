# golden example tests

The tests in this package are used to:

1. that gazelle runs properly (generates correct output from BUILD.in to
   BUILD.out)
2. that the build artifacts can actually be built (bazel build //...).
3. Generate documentation.

`//example/golden:golden_test` runs the gazelle tests.  Each subdirectory of
`testdata/*` represents a separate test.  The `BUILD.in` file is used as the
base, and `BUILD.out` is expected to be generated.  If there is a diff against
the actual `BUILD.out`, the test fails.

The individual `gazelle_testdata_example` tests package up specific examples and
generate a `go_bazel_test` for the example.  So for `//example/golden:java`, a
`//example/golden:java_test` will run a test that actually tries to build
everything within the workspace based on the contents of `BUILD.out`.

# debugging golden example tests

Have any `BUILD.bazel` files crept into the `vendor/` tree?  These are globbed
in `//:all_files` and the presence of one will elide files in that package from
the glob.  You might not see them because they are .gitignored.

```
bazel query //:all_files --output build > all_files.BUILD
```

Try and build from the main test workspace.  To do that, open up rules_go in the
external tree and disable the cleanup function.  For example:

```
$ bazel info output_base
/private/var/tmp/_bazel_pcj/092d6dadaf86f07590903c45033f576e
$ (cd /private/var/tmp/_bazel_pcj/092d6dadaf86f07590903c45033f576e/external/io_bazel_rules_go && code .)
$ code go/tools/bazel_testing/bazel_testing.go
```

```go
	workspaceDir, cleanup, err := setupWorkspace(args, files)
	defer func() {
      < add premature return here to skip cleanup >
		if err := cleanup(); err != nil {
			fmt.Fprintf(os.Stderr, "cleanup error: %v\n", err)
			// Don't fail the test on a cleanup error.
			// Some operating systems (windows, maybe also darwin) can't reliably
			// delete executable files after they're run.
		}
	}()
```

Then, go into the directory that the `go_bazel_test` creates and run bazel there directly:

```sh
$ pushd /private/var/tmp/_bazel_pcj/092d6dadaf86f07590903c45033f576e/bazel_testing/bazel_go_test/main
$ bazel build ...
```
