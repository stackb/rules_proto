module github.com/stackb/rules_proto

go 1.15

require (
	github.com/bazelbuild/bazel-gazelle v0.24.0
	github.com/bazelbuild/buildtools v0.0.0-20221004120235-7186f635531b
	github.com/bazelbuild/rules_go v0.27.0
	github.com/bmatcuk/doublestar v1.2.2
	github.com/emicklei/proto v1.9.0
	github.com/google/go-cmp v0.5.8
	github.com/pmezard/go-difflib v1.0.0
	github.com/stretchr/testify v1.7.0
	go.starlark.net v0.0.0-20220328144851-d1966c6b9fcd
	golang.org/x/sys v0.0.0-20221013171732-95e765b1cc43 // indirect
	gopkg.in/yaml.v3 v3.0.0-20210107192922-496545a6307b
)

// TODO(pcj) Remove once https://github.com/bazelbuild/bazel-gazelle/pull/1033 is merged
replace github.com/bazelbuild/bazel-gazelle => github.com/wolfd/bazel-gazelle v0.0.0-20210917215910-a5bd0e0069da
