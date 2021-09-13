cleanup: go_mod_tidy go_deps buildfiles
	@echo "Done."

tidy:
	bazel run @go_sdk//:bin/go -- mod tidy
	bazel run //:update_go_deps

.PHONY: lint
lint:
	bazel run //:buildifier

.PHONY: gazelle
gazelle:
	bazel run //:gazelle

.PHONY: deps
deps:
	bazel build //deps:*
	cp ./bazel-bin/deps/core_deps.bzl deps/core_deps.bzl
	cp ./bazel-bin/deps/protobuf_deps.bzl deps/protobuf_deps.bzl

.PHONY: site
site:
	bazel build //example/golden:*
	cp -f ./bazel-bin/example/golden/*.md docs/
