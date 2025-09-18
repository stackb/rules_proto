.PHONY: tidy
tidy:
	bazel run @go_sdk//:bin/go -- mod tidy
	bazel run @go_sdk//:bin/go -- mod vendor
	find vendor -name 'BUILD.bazel' | xargs rm
	bazel run //:gazelle
	bazel mod tidy

.PHONY: build
build:
	bazel build //...

.PHONY: test
test:
	bazel test //...

.PHONY: golden_test
golden_test:
	bazel test //example/golden:golden_test --test_output=streamed

.PHONY: site
site:
	bazel build '//example/golden:*'
	cp -f ./bazel-bin/example/golden/*.md docs/

update_pnpm_lock:
	# nvm use 18
	pnpm install --lockfile-only
