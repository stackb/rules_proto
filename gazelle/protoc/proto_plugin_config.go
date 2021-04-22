package protoc

import (
	"github.com/bazelbuild/bazel-gazelle/label"
)

type protoPluginConfig struct {
	Label  label.Label
	Name   string
	Plugin ProtoPlugin
}
