package plugin

import (
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin("grpc:grpc:cpp", &GrpcGrpcCppPlugin{})
}

// GrpcGrpcCppPlugin implements Plugin for the built-in protoc python plugin.
type GrpcGrpcCppPlugin struct{}

// Label implements part of the Plugin interface.
func (p *GrpcGrpcCppPlugin) Label() label.Label {
	return label.New("build_stack_rules_proto", "plugin/grpc/grpc", "cpp")
}

// ShouldApply implements part of the Plugin interface.
func (p *GrpcGrpcCppPlugin) ShouldApply(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface.
func (p *GrpcGrpcCppPlugin) Outputs(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		if !f.HasServices() {
			continue
		}
		name := protoc.PackageFileName(f)
		srcs = append(srcs, name+".grpc.pb.cc", name+".grpc.pb.h")
	}
	return srcs
}
