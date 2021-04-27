package protoc

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const CcProtoName = "cc_proto"

func init() {
	protoc.Plugins().MustRegisterPlugin(CcProtoName, &CcProtoPlugin{})
}

// CcProtoPlugin implements Plugin for the built-in protoc python plugin.
type CcProtoPlugin struct{}

// Label implements part of the Plugin interface.
func (p *CcProtoPlugin) Label() label.Label {
	return label.New("build_stack_rules_proto", "protocolbuffers/protobuf", "cc_plugin")
}

// ShouldApply implements part of the Plugin interface.
func (p *CcProtoPlugin) ShouldApply(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface.
func (p *CcProtoPlugin) Outputs(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		base := f.Name
		pkg := f.Package()
		if pkg.Name != "" {
			base = path.Join(strings.ReplaceAll(pkg.Name, ".", "/"), base)
		}
		if f.HasMessages() || f.HasEnums() {
			srcs = append(srcs, base+".pb.cc", base+".pb.h")
		}
	}
	return srcs
}
