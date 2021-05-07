package golang

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	// ProtocGenGoGrpcName is the name the plugin is registered under.
	ProtocGenGoGrpcName = "protoc-gen-go-grpc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(ProtocGenGoGrpcName, &ProtocGenGoGrpcPlugin{})
}

// ProtocGenGoGrpcPlugin implements Plugin for the the gogo_* family of plugins.
type ProtocGenGoGrpcPlugin struct{}

// Label implements part of the Plugin interface.
func (p *ProtocGenGoGrpcPlugin) Label() label.Label {
	return label.New("build_stack_rules_proto", "grpc/grpc-go", "protoc-gen-go-grpc")
}

func (p *ProtocGenGoGrpcPlugin) ShouldApply(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return true
		}
	}
	return false
}

// Outputs implements part of the Plugin interface
func (p *ProtocGenGoGrpcPlugin) Outputs(rel string, cfg protoc.PackageConfig, lib protoc.ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		if !f.HasServices() {
			continue
		}

		base := f.Name
		pkg := f.Package()
		// see https://github.com/gogo/protobuf/blob/master/protoc-gen-gogo/generator/generator.go#L347
		if goPackage, _, ok := protoc.GoPackageOption(f.Options()); ok {
			base = path.Join(goPackage, base)
		} else if pkg.Name != "" {
			base = path.Join(strings.ReplaceAll(pkg.Name, ".", "/"), base)
		}
		srcs = append(srcs, base+"_grpc.pb.go")
	}
	return srcs
}
