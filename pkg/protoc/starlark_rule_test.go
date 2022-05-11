package protoc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/google/go-cmp/cmp"
)

func TestLoadStarlarkRule(t *testing.T) {
	for name, tc := range map[string]struct {
		code        string
		rc          *LanguageRuleConfig
		pc          *ProtocConfiguration
		wantErr     error
		wantPrinted string
		want        *rule.Rule
	}{
		"degenerate": {
			wantErr: fmt.Errorf(`test.star: rule "test" was never declared`),
		},
		"wrong rule name": {
			code: `
protoc.Rule(
	name = "not-test",
	load_info = lambda: None,
	kind_info = lambda: None,
	provide_rule = lambda rctx, pctx: None,
)
			`,
			wantErr: fmt.Errorf(`test.star: rule "test" was never declared`),
		},
		"missing provide_rule attribute": {
			code: `
protoc.Rule(
	load_info = lambda: None,
	kind_info = lambda: None,
	name = "test", 
)
			`,
			wantErr: fmt.Errorf(`test.star: eval: Rule: missing argument for provide_rule`),
		},
		"provide_rule attribute not callable": {
			code: `
protoc.Rule(
	name = "test", 
	load_info = lambda: None,
	kind_info = lambda: None,
	provide_rule = "not-callable",
)
			`,
			wantErr: fmt.Errorf(`test.star: eval: Rule: for parameter "provide_rule": got string, want callable`),
		},
		"simple": {
			code: `
def make_py_library_rule(self):
	rule = gazelle.Rule(
		name = "py_library",
		kind = "py_library",
	)
	return rule

def provide_rule(rctx, pctx):
	print(rctx)
	print(pctx)
	return struct(
		name = rctx.name,
		kind = "py_library",
		rule = make_py_library_rule,
	)

protoc.Rule(
	name = "test",
	load_info = lambda: None,
	kind_info = lambda: None,
	provide_rule = provide_rule,
)
`,
			rc: &LanguageRuleConfig{
				Config:  &config.Config{},
				Name:    "test",
				Options: map[string]bool{"grpc": true},
			},
			pc: &ProtocConfiguration{},
			wantPrinted: `LanguageRuleConfig(config = Config(repo_name = "", repo_root = "", work_dir = ""), deps = {}, enabled = False, implementation = "", name = "test", options = {"grpc": True}, visibility = {})
ProtocConfiguration(imports = [], language_config = LanguageConfig(enabled = False, name = "", plugins = {}, protoc = "", rules = {}), mappings = {}, outputs = [], package_config = PackageConfig(config = Config(repo_name = "", repo_root = "", work_dir = "")), plugins = [], prefix = "", rel = "")
`,
		},
	} {
		t.Run(name, func(t *testing.T) {
			var err error
			var gotPrinted strings.Builder
			var rule LanguageRule
			rule, err = loadStarlarkLanguageRule("test", "test.star", strings.NewReader(tc.code), func(msg string) {
				gotPrinted.WriteString(msg)
				gotPrinted.Write([]byte{'\n'})
			}, func(loadErr error) {
				err = loadErr
			})
			if err != nil {
				if tc.wantErr != nil {
					if diff := cmp.Diff(tc.wantErr.Error(), err.Error()); diff != "" {
						t.Fatalf("StarlarkRule.load error (-want +got):\n%s", diff)
					}
					return
				} else {
					t.Fatalf("StarlarkRule.load error: %v", err)
				}
			}

			provider := rule.ProvideRule(tc.rc, tc.pc)
			if err != nil {
				if tc.wantErr != nil {
					if diff := cmp.Diff(tc.wantErr.Error(), err.Error()); diff != "" {
						t.Fatalf("StarlarkRule.ProvideRule error (-want +got):\n%s", diff)
					}
					return
				} else {
					t.Fatalf("StarlarkRule.ProvideRule error: %v", err)
				}
			}
			if provider == nil {
				t.Fatalf("StarlarkRule.ProvideRule returned nil")
			}

			t.Log(gotPrinted.String())
			if diff := cmp.Diff(tc.wantPrinted, gotPrinted.String()); diff != "" {
				t.Errorf("StarlarkRule print (-want +got):\n%s", diff)
			}

			got := provider.Rule()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("StarlarkRule.ProvideRule (-want +got):\n%s", diff)
			}
		})
	}
}
