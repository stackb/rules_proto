.PHONY: tidy
tidy: deps
	bazel run @go_sdk//:bin/go -- mod tidy
	bazel run @go_sdk//:bin/go -- mod vendor
	find vendor -name 'BUILD.bazel' | xargs rm
	bazel run //:update_go_deps
	bazel run //:buildifier
	bazel run //:gazelle

.PHONY: deps
deps:
	bzl build //deps:*
	cp -f ./bazel-bin/deps/*.bzl deps/
	chmod 0644 deps/*.bzl
	bzl run //:buildifier -- deps/
		
.PHONY: site
site:
	bazel build //example/golden:*
	cp -f ./bazel-bin/example/golden/*.md docs/

.PHONY: test
test:
	bazel test //example/... //pkg/... //plugin/... //language/... //rules/... //toolchain/... \
		--deleted_packages=//plugin/grpc-ecosystem/grpc-gateway
