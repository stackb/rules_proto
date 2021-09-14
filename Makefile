.PHONY: tidy
tidy: deps
	bazel run @go_sdk//:bin/go -- mod tidy
	bazel run //:update_go_deps
	bazel run //:buildifier
	bazel run //:gazelle

.PHONY: deps
deps:
	bazel build //deps:*
	cp -f ./bazel-bin/deps/core_deps.bzl deps/core_deps.bzl
	cp -f ./bazel-bin/deps/protobuf_core_deps.bzl deps/protobuf_core_deps.bzl
	cp -f ./bazel-bin/deps/grpc_core_deps.bzl deps/grpc_core_deps.bzl
	cp -f ./bazel-bin/deps/grpc_java_deps.bzl deps/grpc_java_deps.bzl

.PHONY: site
site:
	bazel build //example/golden:*
	cp -f ./bazel-bin/example/golden/*.md docs/
