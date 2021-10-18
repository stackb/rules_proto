package protoc

import (
	"testing"

	"github.com/bazelbuild/bazel-gazelle/rule"
)

func TestOtherProtoLibrary(t *testing.T) {
	for name, tc := range map[string]struct {
		rule     *rule.Rule
		wantName string
		wantBase string
		wantSrcs []string
		wantDeps []string
	}{
		"prototypical": {
			rule: withProtoLibraryRule("foo_proto",
				withRuleSrcs("foo.py"),
				withRuleDeps("@a//b:c"),
			),
			wantName: "foo_proto",
			wantBase: "foo",
			wantSrcs: listOf("foo.py"),
			wantDeps: listOf("@a//b:c"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			lib := OtherProtoLibrary{rule: tc.rule}
			name := lib.Name()
			base := lib.BaseName()
			srcs := lib.Srcs()
			deps := lib.Deps()

			if tc.wantName != name {
				t.Errorf("name: want %s, got %s", tc.wantName, name)
			}
			if tc.wantBase != base {
				t.Errorf("base: want %s, got %s", tc.wantBase, base)
			}
			if len(tc.wantSrcs) != len(srcs) {
				t.Fatalf("srcs: want %d, got %d", len(tc.wantSrcs), len(srcs))
			}
			if len(tc.wantDeps) != len(deps) {
				t.Fatalf("deps: want %d, got %d", len(tc.wantDeps), len(deps))
			}
			for i, got := range srcs {
				want := tc.wantSrcs[i]
				if want != got {
					t.Errorf("srcs %d: want %s, got %s", i, want, got)
				}
			}
			for i, got := range deps {
				want := tc.wantDeps[i]
				if want != got {
					t.Errorf("deps %d: want %s, got %s", i, want, got)
				}
			}
		})
	}
}

type ruleOption func(r *rule.Rule)

func withProtoLibraryRule(name string, opts ...ruleOption) *rule.Rule {
	r := rule.NewRule("proto_library", name)
	for _, opt := range opts {
		opt(r)
	}
	return r
}

func withRuleSrcs(srcs ...string) ruleOption {
	return func(r *rule.Rule) {
		r.SetAttr("srcs", srcs)
	}
}

func withRuleDeps(deps ...string) ruleOption {
	return func(r *rule.Rule) {
		r.SetAttr("deps", deps)
	}
}
