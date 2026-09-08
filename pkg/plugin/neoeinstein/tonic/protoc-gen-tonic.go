package tonic

import (
	"path"
	"sort"
	"strings"

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

	perFile := ctx.IsProtoFileMode()
	outputs := p.outputs(ctx.ProtoLibrary, perFile)
	if len(outputs) == 0 {
		return nil
	}

	options := ctx.PluginConfig.GetOptions()
	if perFile {
		// Must match the prost plugin's per_file=true so tonic emits one
		// `<stem>.<pkg>.tonic.rs` per file AND inserts its module entry
		// into the matching `<stem>.<pkg>.rs` rather than the per-package
		// `<pkg>.rs` (which doesn't exist in per-file mode).
		options = append(options, "per_file=true")
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/neoeinstein/tonic", "protoc-gen-tonic"),
		Outputs: outputs,
		Options: options,
	}
}

// ResolvePluginOptions implements the PluginOptionsResolver interface.
// It computes extern_path options based on transitive proto file dependencies
// AND emits self-extern overrides for the library's own packages — needed
// because tonic-generated client/server code references prost types via
// crate-qualified paths and would otherwise be shadowed by parent extern
// crate references through prost's longest-prefix matching.
func (p *ProtocGenTonicPlugin) ResolvePluginOptions(cfg *protoc.PluginConfiguration, r *rule.Rule, from label.Label) []string {
	return prost.ResolveExternPathOptionsForReferences(cfg, r, from)
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

// outputs computes the output files for the plugin.
//
// Per-package mode (default): one `<pkg>.tonic.rs` per proto package that
// declares services.
//
// Per-file mode (the surrounding bazel package opts into `gazelle:proto file`):
// one `<file_stem>.<pkg>.tonic.rs` per .proto file that declares services.
// The vendored protoc-gen-tonic honours this naming convention when
// `per_file=true` is in the options block.
func (p *ProtocGenTonicPlugin) outputs(lib protoc.ProtoLibrary, perFile bool) []string {
	outputs := make([]string, 0)

	if perFile {
		for _, f := range lib.Files() {
			if !f.HasServices() {
				continue
			}
			pkg := f.Package()
			if pkg.Name == "" {
				continue
			}
			stem := strings.TrimSuffix(f.Basename, ".proto")
			filename := stem + "." + pkg.Name + ".tonic.rs"
			if f.Dir != "" {
				filename = path.Join(f.Dir, filename)
			}
			outputs = append(outputs, filename)
		}
		sort.Strings(outputs)
		return outputs
	}

	seen := make(map[string]bool)
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
