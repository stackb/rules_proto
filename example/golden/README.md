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