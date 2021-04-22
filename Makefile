cleanup: go_mod_tidy go_deps buildfiles
	@echo "Done."

go_mod_tidy:
	bazel run @go_sdk//:bin/go -- mod tidy

go_mod_init:
	bazel run @go_sdk//:bin/go -- mod init

go_deps:
	bazel run //:update_go_deps

buildfiles:
	bazel run //:update_build_files

.PHONY: gazelle
gazelle:
	bazel run //:gazelle

gazelle_protoc_test:
	bazel test //gazelle/protoc:protoc_test