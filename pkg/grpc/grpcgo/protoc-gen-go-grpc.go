package grpcgo

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenGoGrpcPlugin{})
}

// ProtocGenGoGrpcPlugin implements Plugin for the the gogo_* family of plugins.
type ProtocGenGoGrpcPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenGoGrpcPlugin) Name() string {
	return "grpc:grpc-go:protoc-gen-go-grpc"
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenGoGrpcPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !p.shouldApply(ctx.ProtoLibrary) {
		return nil
	}
	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/grpc/grpc-go", "protoc-gen-go-grpc"),
		Outputs: p.outputs(ctx.ProtoLibrary),
	}
}

func (p *ProtocGenGoGrpcPlugin) shouldApply(lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return true
		}
	}
	return false
}

func (p *ProtocGenGoGrpcPlugin) outputs(lib protoc.ProtoLibrary) []string {
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
