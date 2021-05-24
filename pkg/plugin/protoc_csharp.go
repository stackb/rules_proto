package plugin

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&ProtocCsharpPlugin{})
}

// ProtocCsharpPlugin implements Plugin for the built-in protoc C# plugin.
type ProtocCsharpPlugin struct{}

// Name implements part of the Plugin interface.
func (p *ProtocCsharpPlugin) Name() string {
	return "protoc:csharp"
}

// Configure implements part of the Plugin interface.
func (p *ProtocCsharpPlugin) Configure(ctx *protoc.PluginContext, cfg *protoc.PluginConfiguration) {
	cfg.Label = label.New("build_stack_rules_proto", "plugin/protoc", "csharp")
	cfg.Outputs = protoc.FlatMapFiles(
		csharpFileName(ctx.Rel, ctx.PluginConfig),
		protoc.Always,
		ctx.ProtoLibrary.Files()...,
	)
	cfg.Out = ctx.Rel
}

func csharpFileName(rel string, cfg protoc.LanguagePluginConfig) func(*protoc.File) []string {
	return func(f *protoc.File) []string {
		// setup the file extension
		ext := ".cs"
		for k, want := range cfg.Options {
			if strings.HasPrefix(k, "file_extension=") && want {
				ext = k[len("file_extension="):]
				continue
			}
		}

		return []string{path.Join(rel, protoc.ToPascalCase(f.Name)) + ext}
	}

}
