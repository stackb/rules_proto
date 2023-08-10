BAZEL := bzl

.PHONY: tidy
tidy: deps
	$(BAZEL) run @go_sdk//:bin/go -- mod tidy
	$(BAZEL) run @go_sdk//:bin/go -- mod vendor
	find vendor -name 'BUILD.bazel' | xargs rm
	$(BAZEL) run //:update_go_deps
	$(BAZEL) run //:buildifier
	$(BAZEL) run //:gazelle

.PHONY: deps
deps:
	$(BAZEL) build //deps:*
	cp -f ./bazel-bin/deps/*.bzl deps/
	chmod 0644 deps/*.bzl
	$(BAZEL) run //:buildifier -- deps/
		
.PHONY: site
site:
	$(BAZEL) build //example/golden:*
	cp -f ./bazel-bin/example/golden/*.md docs/

.PHONY: test
test:
	$(BAZEL) test //example/... //pkg/... //plugin/... //language/... //rules/... //toolchain/... \
		--deleted_packages=//plugin/grpc-ecosystem/grpc-gateway

.PHONY: get
get:
	$(BAZEL) run @go_sdk//:bin/go -- get github.com/bazelbuild/bazel-gazelle@v0.31.0
	$(BAZEL) run @go_sdk//:bin/go -- mod download github.com/bazelbuild/buildtools
	$(BAZEL) run @go_sdk//:bin/go -- mod vendor

update_pnpm_lock:
	# nvm use 18
	pnpm install --lockfile-only
