go_mod_tidy:
	bazel run @go_sdk//:bin/go -- mod tidy

go_mod_init:
	bazel run @go_sdk//:bin/go -- mod init

go_deps:
	bazel run //:update_go_deps

buildfiles:
	bazel run //:gazelle

gazelle_protoc:
	bazel build //gazelle/protoc