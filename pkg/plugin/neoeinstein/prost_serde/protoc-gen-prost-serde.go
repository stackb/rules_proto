package prost_serde

import (
	"path"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/v4/pkg/plugin/neoeinstein/prost"
	"github.com/stackb/rules_proto/v4/pkg/protoc"
)

const ProtocGenProstSerdePluginName = "neoeinstein:prost:protoc-gen-prost-serde"

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocGenProstSerdePlugin{})
}

// ProtocGenProstSerdePlugin implements Plugin for protoc-gen-prost-serde.
type ProtocGenProstSerdePlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocGenProstSerdePlugin) Name() string {
	return ProtocGenProstSerdePluginName
}

// Configure implements part of the Plugin interface.
func (p *ProtocGenProstSerdePlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	if !p.shouldApply(ctx.ProtoLibrary) {
		return nil
	}

	outputs := p.outputs(ctx.ProtoLibrary)
	if len(outputs) == 0 {
		return nil
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/neoeinstein/prost-serde", "protoc-gen-prost-serde"),
		Outputs: outputs,
		Options: ctx.PluginConfig.GetOptions(),
	}
}

// ResolvePluginOptions implements the PluginOptionsResolver interface.
// It computes extern_path options based on transitive proto file dependencies
// AND emits self-extern overrides for the library's own packages — needed
// because prost-serde generates impl blocks at crate-root using absolute
// crate-qualified paths and would otherwise be shadowed by parent extern
// crate references through prost's longest-prefix matching.
func (p *ProtocGenProstSerdePlugin) ResolvePluginOptions(cfg *protoc.PluginConfiguration, r *rule.Rule, from label.Label) []string {
	return prost.ResolveExternPathOptionsForReferences(cfg, r, from)
}

// shouldApply returns true if the library has files with messages or enums.
func (p *ProtocGenProstSerdePlugin) shouldApply(lib protoc.ProtoLibrary) bool {
	for _, f := range lib.Files() {
		if f.HasMessages() || f.HasEnums() {
			return true
		}
	}
	return false
}

// outputs computes the output files for the plugin. Prost-serde generates one
// .serde.rs file per proto package. The path includes the file's directory so
// that mergeSources can handle the rel stripping.
func (p *ProtocGenProstSerdePlugin) outputs(lib protoc.ProtoLibrary) []string {
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

		filename := pkg.Name + ".serde.rs"
		if f.Dir != "" {
			filename = path.Join(f.Dir, filename)
		}
		outputs = append(outputs, filename)
	}

	sort.Strings(outputs)
	return outputs
}
