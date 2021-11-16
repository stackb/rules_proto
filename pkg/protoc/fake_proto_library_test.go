package protoc

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func init() {
	Plugins().MustRegisterPlugin(&fakePlugin{})
	Rules().MustRegisterRule("fake_proto_library", &fakeProtoLibrary{})
}

// fakePlugin implements a mock Plugin
type fakePlugin struct{}

// Name implements part of the Plugin interface.
func (p *fakePlugin) Name() string {
	return "protoc:fake"
}

// Configure implements part of the Plugin interface
func (p *fakePlugin) Configure(ctx *PluginContext) *PluginConfiguration {
	return &PluginConfiguration{
		Label:   label.New("build_stack_rules_proto", "plugin/builtin", "fake"),
		Outputs: p.outputs(ctx.ProtoLibrary),
		Options: p.options(ctx.ProtoLibrary),
	}
}

func (p *fakePlugin) outputs(lib ProtoLibrary) []string {
	srcs := make([]string, 0)
	for _, f := range lib.Files() {
		base := f.Name
		pkg := f.Package()
		if pkg.Name != "" {
			base = path.Join(strings.ReplaceAll(pkg.Name, ".", "/"), base)
		}
		if f.HasMessages() || f.HasEnums() {
			srcs = append(srcs, base+"_fake.pb.go")
		}
	}
	return srcs
}

// Options computes additional options for the plugin.  If the library contains
// services, apply the grpc plugin.
func (p *fakePlugin) options(lib ProtoLibrary) []string {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return []string{"plugins=grpc"}
		}
	}
	return nil
}

// fakeProtoLibrary implements a mock LanguageRule
type fakeProtoLibrary struct{}

// Name implements part of the LanguageRule interface.
func (s *fakeProtoLibrary) Name() string {
	return "fake_proto_library"
}

// KindInfo implements part of the LanguageRule interface.
func (s *fakeProtoLibrary) KindInfo() rule.KindInfo {
	return rule.KindInfo{}
}

// LoadInfo implements part of the LanguageRule interface.
func (s *fakeProtoLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *fakeProtoLibrary) ProvideRule(rc *LanguageRuleConfig, pc *ProtocConfiguration) RuleProvider {
	return nil
}
