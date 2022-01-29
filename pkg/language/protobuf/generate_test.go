package protobuf

import (
	"os"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/bazel-gazelle/testtools"
	"github.com/google/go-cmp/cmp"
)

func TestGenerateRules(t *testing.T) {
	for name, tc := range map[string]struct {
		rel       string
		files     []testtools.FileSpec
		args      language.GenerateArgs
		want      language.GenerateResult
		pre, post func(*testGenerateRulesState)
	}{
		"empty case": {
			args: language.GenerateArgs{
				Config: makeTestConfig(""),
			},
			want: language.GenerateResult{
				Gen:     []*rule.Rule{},
				Empty:   []*rule.Rule{},
				Imports: []interface{}{},
			},
		},
		"registers labels qualified with the config.RepoName": {
			files: []testtools.FileSpec{
				{
					Path: "foo.proto",
					Content: `syntax = "proto3";
import "google/protobuf/any.proto";
					`,
				},
			},
			args: language.GenerateArgs{
				Config:       makeTestConfig("contoso"),
				RegularFiles: []string{"foo.proto"},
				OtherGen:     []*rule.Rule{makeTestProtoLibraryRule()},
			},
			want: language.GenerateResult{
				Gen:     []*rule.Rule{},
				Empty:   []*rule.Rule{},
				Imports: []interface{}{},
			},
			post: func(state *testGenerateRulesState) {
				wantProvided := []importResolverProvide{
					// this says "foo.proto has a dependency on any.proto".  The label isn't currently used.
					{
						lang:    "proto",
						impLang: "depends",
						imp:     "foo.proto",
						label:   label.New("", "google/protobuf", "any.proto"),
					},
					// this informs the resolver that "messages.proto" is provided by the proto_library @contoso//:foo_library
					{
						lang:    "proto",
						impLang: "proto",
						imp:     "messages.proto",
						label:   label.New("contoso", "", "foo_library"),
					},
				}
				if diff := cmp.Diff(wantProvided, state.resolver.provided, cmp.AllowUnexported(importResolverProvide{})); diff != "" {
					t.Error("unexpected diff:", diff)
				}
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			dir, cleanup := testtools.CreateFiles(t, tc.files)
			defer cleanup()

			if err := os.Chdir(dir); err != nil {
				t.Fatal(err)
			}

			ext := NewProtobufLang("test")
			tc.args.Rel = tc.rel
			tc.args.Config.WorkDir = dir

			state := &testGenerateRulesState{
				t:        t,
				tmpdir:   dir,
				ext:      ext,
				resolver: &mockImportResolver{},
			}
			ext.resolver = state.resolver

			if tc.pre != nil {
				tc.pre(state)
			}

			got := ext.GenerateRules(tc.args)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Error("unexpected diff:", diff)
			}

			if tc.post != nil {
				tc.post(state)
			}
		})
	}
}

type testGenerateRulesState struct {
	t        *testing.T
	tmpdir   string
	ext      *protobufLang
	resolver *mockImportResolver
}

func makeTestProtoLibraryRule() *rule.Rule {
	r := rule.NewRule("proto_library", "foo_library")
	r.SetAttr("srcs", []string{"messages.proto"})
	return r
}

func makeTestConfig(repoName string) *config.Config {
	return &config.Config{
		RepoName: repoName,
		Exts:     make(map[string]interface{}),
	}
}

type importResolverProvide struct {
	lang, impLang, imp string
	label              label.Label
}

type mockImportResolver struct {
	provided []importResolverProvide
	resolved []resolve.FindResult
}

func (m *mockImportResolver) Resolve(lang, impLang, imp string) []resolve.FindResult {
	return m.resolved
}

func (m *mockImportResolver) Provide(lang string, impLang, val string, location label.Label) {
	m.provided = append(m.provided, importResolverProvide{lang, impLang, val, location})
}
