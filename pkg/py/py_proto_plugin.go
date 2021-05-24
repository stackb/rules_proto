package protoc

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocPythonPlugin{})
}

// ProtocPythonPlugin implements Plugin for the built-in protoc python plugin.
type ProtocPythonPlugin struct{}

// Label implements part of the Plugin interface.
func (p *ProtocPythonPlugin) Name() string {
	return "protoc:python"
}

// Label implements part of the Plugin interface.
func (p *ProtocPythonPlugin) Label() label.Label {
	return label.New("build_stack_rules_proto", "protocolbuffers/protobuf", "py_proto_plugin")
}

// ShouldApply implements part of the Plugin interface.
func (p *ProtocPythonPlugin) ShouldApply(ctx *protoc.PluginContext) bool {
	for _, f := range ctx.ProtoLibrary.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface.
func (p *ProtocPythonPlugin) Outputs(ctx *protoc.PluginContext) []string {
	srcs := make([]string, 0)
	for _, f := range ctx.ProtoLibrary.Files() {
		base := f.Name
		pkg := f.Package()
		if pkg.Name != "" {
			base = path.Join(strings.ReplaceAll(pkg.Name, ".", "/"), base)
		}
		if f.HasMessages() || f.HasEnums() {
			srcs = append(srcs, base+"_pb2.py")
		}
	}
	return srcs
}
