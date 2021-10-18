module github.com/stackb/rules_proto

go 1.15

require (
	github.com/bazelbuild/bazel-gazelle v0.24.0
	github.com/bazelbuild/buildtools v0.0.0-20210408102303-2b0a1af1a898
	github.com/bazelbuild/rules_go v0.27.0
	github.com/emicklei/proto v1.9.0
	github.com/google/go-cmp v0.5.5
	github.com/pmezard/go-difflib v1.0.0
	github.com/stretchr/testify v1.7.0
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

// TODO(pcj) Remove once https://github.com/bazelbuild/bazel-gazelle/pull/1033 is merged
replace github.com/bazelbuild/bazel-gazelle => github.com/wolfd/bazel-gazelle v0.0.0-20210917215910-a5bd0e0069da
