package protoc

import (
	"fmt"
	"strings"
	"testing"

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
	} {
		t.Run(name, func(t *testing.T) {
			var err error
			var gotPrinted strings.Builder
			var rule LanguageRule
			rule, err = loadStarlarkLanguageRule("test", "test.star", strings.NewReader(tc.code), func(msg string) {
				gotPrinted.WriteString(msg)
				gotPrinted.Write([]byte{'\n'})
			}, func(configureErr error) {
				err = configureErr
			})
			if err != nil {
				if tc.wantErr != nil {
					if diff := cmp.Diff(tc.wantErr.Error(), err.Error()); diff != "" {
						t.Fatalf("StarlarkRule.Configure error (-want +got):\n%s", diff)
					}
					return
				} else {
					t.Fatalf("StarlarkRule.Configure error: %v", err)
				}
			}

			t.Log(gotPrinted.String())
			if diff := cmp.Diff(tc.wantPrinted, gotPrinted.String()); diff != "" {
				t.Errorf("StarlarkRule.Configure print (-want +got):\n%s", diff)
			}

			provider := rule.ProvideRule(tc.rc, tc.pc)
			if provider == nil {
				t.Fatalf("StarlarkRule.ProvideRule returned nil")
			}

			got := provider.Rule()
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("StarlarkRule.Configure (-want +got):\n%s", diff)
			}
		})
	}
}
