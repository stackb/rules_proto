package protoc

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const CcGrpcName = "cc_grpc"

func init() {
	protoc.Plugins().MustRegisterPlugin(CcGrpcName, &CcGrpcPlugin{})
}

// CcGrpcPlugin implements Plugin for the built-in protoc python plugin.
type CcGrpcPlugin struct{}

// Label implements part of the Plugin interface.
func (p *CcGrpcPlugin) Label() label.Label {
	return label.New("build_stack_rules_proto", "grpc/grpc", "cc_plugin")
}

// ShouldApply implements part of the Plugin interface.
func (p *CcGrpcPlugin) ShouldApply(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface.
func (p *CcGrpcPlugin) Outputs(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		base := f.Name
		pkg := f.Package()
		if pkg.Name != "" {
			base = path.Join(strings.ReplaceAll(pkg.Name, ".", "/"), base)
		}
		if f.HasServices() {
			srcs = append(srcs, base+".grpc.pb.cc", base+".grpc.pb.h")
		}
	}
	return srcs
}
