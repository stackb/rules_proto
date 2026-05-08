package tonic

import (
	"path"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/v4/pkg/plugin/neoeinstein/prost"
	"github.com/stackb/rules_proto/v4/pkg/protoc"
)

const ProtocGenTonicPluginName = "neoeinstein:prost:protoc-gen-tonic"

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenTonicPlugin{})
}

// ProtocGenTonicPlugin implements Plugin for protoc-gen-tonic.
type ProtocGenTonicPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenTonicPlugin) Name() string {
	return ProtocGenTonicPluginName
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenTonicPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !p.shouldApply(ctx.ProtoLibrary) {
		return nil
	}

	outputs := p.outputs(ctx.ProtoLibrary)
	if len(outputs) == 0 {
		return nil
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/neoeinstein/tonic", "protoc-gen-tonic"),
		Outputs: outputs,
		Options: ctx.PluginConfig.GetOptions(),
	}
}

// ResolvePluginOptions implements the PluginOptionsResolver interface.
// It computes extern_path options based on transitive proto file dependencies.
func (p *ProtocGenTonicPlugin) ResolvePluginOptions(cfg *protoc.PluginConfiguration, r *rule.Rule, from label.Label) []string {
	return prost.ResolveExternPathOptions(cfg, r, from)
}

// shouldApply returns true if the library has files with services.
func (p *ProtocGenTonicPlugin) shouldApply(lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return true
		}
	}
	return false
}

// outputs computes the output files for the plugin. Tonic generates one
// .tonic.rs file per proto package that has services. The path includes the
// file's directory so that mergeSources can handle the rel stripping.
func (p *ProtocGenTonicPlugin) outputs(lib protoc.ProtoLibrary) []string {
	seen := make(map[string]bool)
	outputs := make([]string, 0)

	for _, f := range lib.Files() {
		if !f.HasServices() {
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

		filename := pkg.Name + ".tonic.rs"
		if f.Dir != "" {
			filename = path.Join(f.Dir, filename)
		}
		outputs = append(outputs, filename)
	}

	sort.Strings(outputs)
	return outputs
}
