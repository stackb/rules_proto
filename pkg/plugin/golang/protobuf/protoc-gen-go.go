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
	mappings := GetImportMappings(ctx.PluginConfig.GetOptions())
	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/golang/protobuf", "protoc-gen-go"),
		Outputs: p.outputs(ctx.ProtoLibrary, mappings),
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

func (p *ProtocGenGoPlugin) outputs(lib protoc.ProtoLibrary, importMappings map[string]string) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		if !(f.HasMessages() || f.HasEnums()) {
			continue
		}
		srcs = append(srcs, GetGoOutputBaseName(f, importMappings)+".pb.go")
	}
	return srcs
}

func GetGoOutputBaseName(f *protoc.File, importMappings map[string]string) string {
	base := f.Name
	pkg := f.Package()
	// see https://github.com/gogo/protobuf/blob/master/protoc-gen-gogo/generator/generator.go#L347
	if mapping := importMappings[path.Join(f.Dir, f.Basename)]; mapping != "" {
		base = path.Join(mapping, base)
	} else if goPackage, _, ok := protoc.GoPackageOption(f.Options()); ok {
		base = path.Join(goPackage, base)
	} else if pkg.Name != "" {
		base = path.Join(strings.ReplaceAll(pkg.Name, ".", "/"), base)
	}
	return base
}

func GetImportMappings(options []string) map[string]string {
	// gather options that look like protoc-gen-go "importmapping" (M) options
	// (e.g Mfoo.proto=github.com/example/foo).
	mappings := make(map[string]string)

	for _, opt := range options {
		if !strings.HasPrefix(opt, "M") {
			continue
		}
		parts := strings.SplitN(opt[1:], "=", 2)
		if len(parts) != 2 {
			continue
		}
		mappings[parts[0]] = parts[1]
	}

	return mappings
}
