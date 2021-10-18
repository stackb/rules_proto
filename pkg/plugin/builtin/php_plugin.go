package builtin

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&PhpPlugin{})
}

// PhpPlugin implements Plugin for the built-in protoc php plugin.
type PhpPlugin struct{}

// Name implements part of the Plugin interface.
func (p *PhpPlugin) Name() string {
	return "builtin:php"
}

// Configure implements part of the Plugin interface.
func (p *PhpPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	return &protoc.PluginConfiguration{
		Label: label.New("build_stack_rules_proto", "plugin/builtin", "php"),
		Outputs: protoc.FlatMapFiles(
			phpFileName(ctx.Rel),
			protoc.Always,
			ctx.ProtoLibrary.Files()...,
		),
		Out:     ctx.Rel,
		Options: ctx.PluginConfig.GetOptions(),
	}
}

func phpFileName(rel string) func(f *protoc.File) []string {
	relDir := strings.Title(rel)

	return func(f *protoc.File) []string {
		outs := make([]string, 0)

		// Compute the base dir where files are generated
		dir := ""
		pkg := f.Package()
		if pkg.Name != "" {
			dir = path.Join(strings.Title(strings.ReplaceAll(pkg.Name, ".", "/")), dir)
		}

		// php_namespace overrides package
		ns, _ := protoc.GetNamedOption(f.Options(), "php_namespace")
		if ns != "" {
			dir = ns
		}

		// Add the metadata file
		mns, _ := protoc.GetNamedOption(f.Options(), "php_metadata_namespace")
		if mns == "" {
			mns = "GPBMetadata"
		}
		outs = append(outs, path.Join(rel, mns, relDir, strings.Title(f.Name))+".php")

		// Add enums
		for _, e := range f.Enums() {
			outs = append(outs, path.Join(dir, rel, strings.Title(e.Name))+".php")
		}

		// Add messages
		for _, m := range f.Messages() {
			outs = append(outs, path.Join(dir, rel, strings.Title(m.Name))+".php")
		}

		return outs
	}
}
