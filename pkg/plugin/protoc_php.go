package plugin

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocPhpPlugin{})
}

// ProtocPhpPlugin implements Plugin for the built-in protoc C++ plugin.
type ProtocPhpPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocPhpPlugin) Name() string {
	return "protoc:php"
}

// Configure implements part of the Plugin interface.
func (p *ProtocPhpPlugin) Configure(ctx *protoc.PluginContext, cfg *protoc.PluginConfiguration) {
	cfg.Label = label.New("build_stack_rules_proto", "plugin/protoc", "php")
	cfg.Outputs = protoc.FlatMapFiles(
		phpFileName,
		protoc.Always,
		ctx.ProtoLibrary.Files()...,
	)
}

func phpFileName(f *protoc.File) []string {
	outs := make([]string, 0)

	// Compute the base dir where files are generated
	dir := ""
	pkg := f.Package()
	if pkg.Name != "" {
		dir = path.Join(strings.Title(strings.ReplaceAll(pkg.Name, ".", "/")), dir)
	}

	// php_namespace overrides package
	ns := protoc.GetNamedOption(f.Options(), "php_namespace")
	if ns != "" {
		dir = ns
	}

	// Add the metadata file
	mns := protoc.GetNamedOption(f.Options(), "php_metadata_namespace")
	if mns == "" {
		mns = "GPBMetadata"
	}
	outs = append(outs, path.Join(mns, strings.Title(f.Name))+".php")

	// Add enums
	for _, e := range f.Enums() {
		outs = append(outs, path.Join(dir, strings.Title(e.Name))+".php")
	}

	// Add messages
	for _, m := range f.Messages() {
		outs = append(outs, path.Join(dir, strings.Title(m.Name))+".php")
	}

	return outs
}
