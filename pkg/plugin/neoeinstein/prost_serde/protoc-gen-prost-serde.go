package prost_serde

import (
	"path"
	"sort"
	"strings"

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

	perFile := ctx.IsProtoFileMode()
	outputs := p.outputs(ctx.ProtoLibrary, perFile)
	if len(outputs) == 0 {
		return nil
	}

	options := ctx.PluginConfig.GetOptions()
	if perFile {
		// Must match the prost plugin's per_file=true so pbjson emits one
		// `<stem>.<pkg>.serde.rs` per file AND points its `insert_include`
		// pragma at the matching `<stem>.<pkg>.rs` (which prost wrote);
		// the unpatched form inserts into `<pkg>.rs`, a filename that no
		// longer exists in per-file mode.
		options = append(options, "per_file=true")
	}

	return &protoc.PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/neoeinstein/prost-serde", "protoc-gen-prost-serde"),
		Outputs: outputs,
		Options: options,
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

// outputs computes the output files for the plugin.
//
// Per-package mode (default): one `<pkg>.serde.rs` per proto package.
//
// Per-file mode (the surrounding bazel package opts into `gazelle:proto file`):
// one `<file_stem>.<pkg>.serde.rs` per .proto file that defines messages or
// enums. The vendored protoc-gen-prost-serde honours this naming convention
// when `per_file=true` is in the options block.
func (p *ProtocGenProstSerdePlugin) outputs(lib protoc.ProtoLibrary, perFile bool) []string {
	outputs := make([]string, 0)

	if perFile {
		for _, f := range lib.Files() {
			if !prost.HasGeneratableContent(f) {
				// Same rationale as the per-package branch below: pbjson
				// emits no file for extend-only protos, and a declared
				// output would trip proto_compile.bzl's `mv` rename.
				continue
			}
			pkg := f.Package()
			if pkg.Name == "" {
				continue
			}
			stem := strings.TrimSuffix(f.Basename, ".proto")
			filename := stem + "." + pkg.Name + ".serde.rs"
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
		if !prost.HasGeneratableContent(f) {
			// Skip files with no serializable types. pbjson-build emits
			// nothing for extend-only files, and the declared output would
			// later fail proto_compile.bzl's `mv` rename.
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
