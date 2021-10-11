// gencopy is a utility program that copies bazel outputs back into the
// workspace source tree.  Ideally, you don't have any generated files committed
// to VCS, but sometimes you do.
//
package rules_go

import (
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/stackb/rules_proto/pkg/protoc"
)

func TestGoLibraryRuleImportPath(t *testing.T) {
	for name, tc := range map[string]struct {
		files       []*protoc.File
		ruleOptions []string
		plugins     []*protoc.PluginConfiguration
		want        string // importpath string
	}{
		"degenerate": {},
		"from go_package option": {
			files: []*protoc.File{
				newProtoFile(t, "proto", "foo.proto",
					`syntax = "proto3";`,
					`option go_package = "github.com/example.com/foo";`,
				),
			},
			want: "github.com/example.com/foo",
		},
		"from rule importmapping option": {
			files: []*protoc.File{
				newProtoFile(t, "proto", "foo.proto"),
			},
			ruleOptions: []string{
				"Mproto/foo.proto=github.com/example.com/foo",
			},
			want: "github.com/example.com/foo",
		},
		"from plugin importmapping option": {
			files: []*protoc.File{
				newProtoFile(t, "proto", "foo.proto"),
			},
			plugins: []*protoc.PluginConfiguration{
				{
					Options: []string{
						"Mproto/foo.proto=github.com/example.com/foo",
					},
				},
			},
			want: "github.com/example.com/foo",
		},
	} {
		t.Run(name, func(t *testing.T) {
			ruleConfig := protoc.NewLanguageRuleConfig(nil, "proto_go_library")
			for _, d := range tc.ruleOptions {
				ruleConfig.Options[d] = true
			}
			gazelleRule := rule.NewRule("proto_library", "foo_proto")
			goRule := &goLibraryRule{
				config: &protoc.ProtocConfiguration{
					Plugins: tc.plugins,
					Library: protoc.NewOtherProtoLibrary(nil, gazelleRule, tc.files...),
				},
				ruleConfig: ruleConfig,
			}
			got := goRule.importPath()
			if tc.want != got {
				t.Errorf("want %q, got %q", tc.want, got)
			}
		})
	}
}

func newProtoFile(t *testing.T, dir, name string, lines ...string) *protoc.File {
	f := protoc.NewFile(dir, name)
	content := strings.Join(lines, "\n")
	in := strings.NewReader(content)
	if err := f.ParseReader(in); err != nil {
		t.Fatal(err)
	}
	return f
}
