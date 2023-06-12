# proto_repository

This test demonstrates the use of proto_repository and go language generation.
A few notes about the test:

- The files are used in both "//example/golden:golden_test" (a golden test) and
  "//example/golden:proto_repository_test" (go_bazel_test).
- The proto mode is "file", so the proto_library rule for "app.proto" has two
  dependencies: one for annotations.proto and another for field_behavior.proto.
  Yet, the proto_go_library has only a single dependency on
  annotations_go_proto.  This is because go sources that share the same
  "importpath" must be compiled together under a single go_library rule.
  Therefore, there is logic that merges proto_go_library rules when they share
  the same importpath.  The relationship of annotations.proto and
  field_behavior.proto tests this merge behavior.

## //example/golden:golden_test (golden test)

- Only the ".gazelle.args and "imports.csv" file is used in the golden test.
  This simulates the existence of the "@googleapis" repo and the imports.csv
  file file that it would produce.  Only the two relevant entries of the file
  are included for deps resolution of the proto_library and proto_go_library
  rules.

## //example/golden:proto_repository_test" (go_bazel_test)

- This runs "bazel build ..." and then "bazel run //:gazelle", "bazel build
  //...".
- The effective WORKSPACE file is a concatenation of "prebuilt.WORKSPACE" and
  the "WORKSPACE" file here.
- The order of loads in that effective workspace are important, namely we need
  to load newer versions of @org_golang_google_grpc before rules_go or gazelle
  try and do that.
- proto file mode is used.
- reresolve_known_proto_imports = True is needed to rewrite labels that would normally
  be in the "go_googleapis" and "com_google_protobuf" external workspace.
