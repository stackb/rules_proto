package prost

import (
	"path"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/v4/pkg/protoc"
)

const (
	ProtocGenProstPluginName = "neoeinstein:prost:protoc-gen-prost"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenProstPlugin{})
}

// ProtocGenProstPlugin implements Plugin for protoc-gen-prost.
type ProtocGenProstPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenProstPlugin) Name() string {
	return ProtocGenProstPluginName
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenProstPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !p.shouldApply(ctx.ProtoLibrary) {
		return nil
	}

	outputs := p.outputs(ctx.ProtoLibrary)
	if len(outputs) == 0 {
		return nil
	}

	p.registerExternPaths(ctx.ProtoLibrary)

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/neoeinstein/prost", "protoc-gen-prost"),
		Outputs: outputs,
		Options: ctx.PluginConfig.GetOptions(),
	}
}

// ResolvePluginOptions implements the PluginOptionsResolver interface.
// It computes extern_path options based on transitive proto file dependencies.
func (p *ProtocGenProstPlugin) ResolvePluginOptions(cfg *protoc.PluginConfiguration, r *rule.Rule, from label.Label) []string {
	return ResolveExternPathOptions(cfg, r, from)
}

// shouldApply returns true if the library has files with messages or enums.
func (p *ProtocGenProstPlugin) shouldApply(lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// outputs computes the output files for the plugin. Prost generates one .rs
// file per proto package, named {proto_package}.rs. The path includes the
// file's directory so that mergeSources can handle the rel stripping.
func (p *ProtocGenProstPlugin) outputs(lib protoc.ProtoLibrary) []string {
	seen := make(map[string]bool)
	outputs := make([]string, 0)

	for _, f := range lib.Files() {
		if !(f.HasMessages() || f.HasEnums()) {
			continue
		}
		pkg := f.Package()
		if pkg.Name == "" {
			continue
		}
		if seen[pkg.Name] {
			continue
		}
		seen[pkg.Name] = true

		filename := pkg.Name + ".rs"
		if f.Dir != "" {
			filename = path.Join(f.Dir, filename)
		}
		outputs = append(outputs, filename)
	}

	sort.Strings(outputs)
	return outputs
}

// registerExternPaths records prost extern_path data in the GlobalResolver for
// each proto file in the library. This data is later consumed by
// ResolveTransitiveExternPaths when computing extern_path options for dependent
// packages.
//
// The label encodes: Pkg = proto package name, Name = crate name.
func (p *ProtocGenProstPlugin) registerExternPaths(lib protoc.ProtoLibrary) {
	crateName := lib.BaseName() + "_rs"

	for _, f := range lib.Files() {
		pkg := f.Package()
		if pkg.Name == "" {
			continue
		}

		protoFile := path.Join(f.Dir, f.Basename)
		protoc.GlobalResolver().Provide(
			"proto",
			"prost_extern",
			protoFile,
			label.New("", pkg.Name, crateName),
		)
	}
}
