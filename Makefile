
.PHONY: tidy
tidy: deps
	bazel run @go_sdk//:bin/go -- mod tidy
	bazel run @go_sdk//:bin/go -- mod vendor
	find vendor -name 'BUILD.bazel' | xargs rm
	bazel run //:update_go_deps
	bazel run //:buildifier
	bazel run //:gazelle

.PHONY: gazelle
gazelle:
	bazel run //:gazelle

.PHONY: deps
deps:
	bazel build //deps:*
	cp -f ./bazel-bin/deps/*.bzl deps/
	chmod 0644 deps/*.bzl
	bazel run //:buildifier -- deps/
		
.PHONY: site
site:
	bazel build //example/golden:*
	cp -f ./bazel-bin/example/golden/*.md docs/

.PHONY: golden_test
golden_test:
	bazel test //example/golden:golden_test --test_output=streamed

.PHONY: example_test
example_test:
	bazel test //example/golden:proto_compiled_sources_test --test_output=streamed

.PHONY: test
test:
	bazel test --keep_going //example/... //pkg/... //plugin/... //language/... //rules/... //toolchain/...

.PHONY: get
get:
	bazel run @go_sdk//:bin/go -- get github.com/bazelbuild/bazel-gazelle@v0.31.0
	bazel run @go_sdk//:bin/go -- mod download github.com/bazelbuild/buildtools
	bazel run @go_sdk//:bin/go -- mod vendor

update_pnpm_lock:
	# nvm use 18
	pnpm install --lockfile-only
