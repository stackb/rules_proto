package rules_python

import (
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/google/go-cmp/cmp"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func makeTestProtoLibrary(cfg ...func(*rule.Rule)) *rule.Rule {
	r := rule.NewRule("proto_library", "test_proto")
	for _, fn := range cfg {
		fn(r)
	}
	return r
}

func TestProtoPyLibraryRule(t *testing.T) {
	for name, tc := range map[string]struct {
		cfg  protoc.LanguageRuleConfig
		pc   protoc.ProtocConfiguration
		want string
	}{
		"degenerate": {
			cfg: *protoc.NewLanguageRuleConfig(config.New(), "py"),
			pc: protoc.ProtocConfiguration{
				Library: protoc.NewOtherProtoLibrary(nil, makeTestProtoLibrary()),
			},
		},
		"simple": {
			cfg: *protoc.NewLanguageRuleConfig(config.New(), "py"),
			pc: protoc.ProtocConfiguration{
				Library: protoc.NewOtherProtoLibrary(nil, makeTestProtoLibrary()),
				Plugins: []*protoc.PluginConfiguration{
					{
						Config: &protoc.LanguagePluginConfig{
							Implementation: "builtin:python",
						},
						Outputs: []string{"test_pb2.py"},
					},
				},
			},
			want: `
proto_py_library(
    name = "test_py_library",
    srcs = ["test_pb2.py"],
)
`,
		},
		"strip_import_prefix": {
			cfg: *protoc.NewLanguageRuleConfig(config.New(), "py"),
			pc: protoc.ProtocConfiguration{
				Library: protoc.NewOtherProtoLibrary(nil, makeTestProtoLibrary(func(r *rule.Rule) {
					r.SetAttr("strip_import_prefix", "/com/foo/")
				})),
				Plugins: []*protoc.PluginConfiguration{
					{
						Config: &protoc.LanguagePluginConfig{
							Implementation: "builtin:python",
						},
						Outputs: []string{"test_pb2.py"},
					},
				},
			},
			want: `
proto_py_library(
    name = "test_py_library",
    srcs = ["test_pb2.py"],
    imports = ["../.."],
)
`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			lib := protoPyLibrary{}
			impl := lib.ProvideRule(&tc.cfg, &tc.pc)
			var got string
			if impl != nil {
				rule := impl.Rule()
				got = printRules(rule)
			}
			if diff := cmp.Diff(strings.TrimSpace(tc.want), strings.TrimSpace(got)); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}

func printRules(rules ...*rule.Rule) string {
	file := rule.EmptyFile("", "")
	for _, r := range rules {
		r.Insert(file)
	}
	return string(file.Format())
}
