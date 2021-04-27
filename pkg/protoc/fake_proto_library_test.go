package protoc

import (
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/rule"
)

func init() {
	Plugins().MustRegisterPlugin("fake_proto", &fakePlugin{})
	Rules().MustRegisterRule("fake_proto_library", &fakeProtoLibrary{})
}

// fakePlugin implements a mock Plugin
type fakePlugin struct{}

func (p *fakePlugin) ShouldApply(rel string, cfg PackageConfig, lib ProtoLibrary) bool {
	return true
}

// Outputs implements part of the Plugin interface
func (p *fakePlugin) Outputs(rel string, cfg PackageConfig, lib ProtoLibrary) []string {
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

// Options implements part of the optional PluginOptionsProvider interface.  If
// the library contains services, apply the grpc plugin.
func (p *fakePlugin) Options(rel string, c *PackageConfig, lib ProtoLibrary) []string {
	for _, f := range lib.Files() {
		if f.HasServices() {
			return []string{"plugins=grpc"}
		}
	}
	return nil
}

// fakeProtoLibrary implements a mock LanguageRule
type fakeProtoLibrary struct{}

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
