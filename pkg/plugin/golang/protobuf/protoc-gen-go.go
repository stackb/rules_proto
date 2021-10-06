package protobuf

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

const ProtocGenGoPluginName = "golang:protobuf:protoc-gen-go"

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenGoPlugin{})
}

// ProtocGenGoPlugin implements Plugin for the the gogo_* family of plugins.
type ProtocGenGoPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenGoPlugin) Name() string {
	return ProtocGenGoPluginName
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenGoPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !p.shouldApply(ctx.ProtoLibrary) {
		return nil
	}
	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/golang/protobuf", "protoc-gen-go"),
		Outputs: p.outputs(ctx.ProtoLibrary),
	}
}

func (p *ProtocGenGoPlugin) shouldApply(lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

func (p *ProtocGenGoPlugin) outputs(lib protoc.ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		if !(f.HasMessages() || f.HasEnums()) {
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
		srcs = append(srcs, base+".pb.go")
	}
	return srcs
}
