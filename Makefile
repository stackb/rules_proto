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
	cp -f ./bazel-bin/deps/prebuilt_protoc_deps.bzl deps/prebuilt_protoc_deps.bzl
	cp -f ./bazel-bin/deps/scala_deps.bzl deps/scala_deps.bzl
	cp -f ./bazel-bin/deps/go_core_deps.bzl deps/go_core_deps.bzl
	cp -f ./bazel-bin/deps/nodejs_deps.bzl deps/nodejs_deps.bzl
	cp -f ./bazel-bin/deps/ts_proto_deps.bzl deps/ts_proto_deps.bzl
	cp -f ./bazel-bin/deps/closure_deps.bzl deps/closure_deps.bzl
	chmod 0644 deps/*.bzl rules/nodejs/deps.bzl
	bazel run //:buildifier -- deps/
		
.PHONY: site
site:
	bazel build //example/golden:*
	cp -f ./bazel-bin/example/golden/*.md docs/
