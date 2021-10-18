package builtin

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func init() {
	protoc.Plugins().MustRegisterPlugin(&CsharpPlugin{})
}

// CsharpPlugin implements Plugin for the built-in protoc C# plugin.
type CsharpPlugin struct{}

// Name implements part of the Plugin interface.
func (p *CsharpPlugin) Name() string {
	return "builtin:csharp"
}

// Configure implements part of the Plugin interface.
func (p *CsharpPlugin) Configure(ctx *protoc.PluginContext) *protoc.PluginConfiguration {
	return &protoc.PluginConfiguration{
		Label: label.New("build_stack_rules_proto", "plugin/builtin", "csharp"),
		Outputs: protoc.FlatMapFiles(
			csharpFileName(ctx.Rel, ctx.PluginConfig),
			protoc.Always,
			ctx.ProtoLibrary.Files()...,
		),
		Out:     ctx.Rel,
		Options: ctx.PluginConfig.GetOptions(),
	}
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
