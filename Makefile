
.PHONY: tidy
tidy: deps
	bazel run @go_sdk//:bin/go -- mod tidy
	bazel run @go_sdk//:bin/go -- mod vendor
	find vendor -name 'BUILD.bazel' | xargs rm
	bazel run //:gazelle

.PHONY: golden_test
golden_test:
	bazel test //example/golden:golden_test --test_output=streamed

.PHONY: test
test:
	bazel test ...

update_pnpm_lock:
	# nvm use 18
	pnpm install --lockfile-only
