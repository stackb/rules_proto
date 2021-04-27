package java

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const JavaGrpcName = "java_grpc"

func init() {
	protoc.Plugins().MustRegisterPlugin(JavaGrpcName, &JavaGrpcPlugin{})
}

// JavaGrpcPlugin implements Plugin for the built-in protoc java plugin.
type JavaGrpcPlugin struct{}

// ShouldApply implements part of the Plugin interface.
func (p *JavaGrpcPlugin) ShouldApply(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return true
		}
	}
	return false
}

// Label implements part of the Plugin interface.
func (p *JavaGrpcPlugin) Label() label.Label {
	return label.New("build_stack_rules_proto", "grpc/grpc-java", "grpc_plugin")
}

// Outputs implements part of the Plugin interface.
func (p *JavaGrpcPlugin) Outputs(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
	return []string{srcjarFile(rel, lib.BaseName()+"_grpc")}
}

// Out implements part the optional PluginOutProvider interface.
func (p *JavaGrpcPlugin) Out(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) string {
	return srcjarFile(rel, lib.BaseName()+"_grpc")
}
