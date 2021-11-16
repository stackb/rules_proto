package protoc

import (
	"testing"

	"github.com/bazelbuild/bazel-gazelle/rule"
)

type packageConfigCheck func(t *testing.T, cfg *PackageConfig)

type packageConfigTestCase struct {
	rel        string
	directives []rule.Directive
	check      packageConfigCheck
	err        error
}

func TestPackageConfigClone(t *testing.T) {
	initialState := allPackageChecks(
		withLanguage("fake",
			withLanguageEnabled(true),
		),
		withLanguageRule("fake_proto_library",
			withLanguageRuleEnabled(true),
		),
	)
	finalState := allPackageChecks(
		withLanguage("fake",
			withLanguageEnabled(true),
		),
		withLanguageRule("fake_proto_library",
			withLanguageRuleEnabled(false),
		),
		withPlugin("fake_proto"), // withPluginToolEquals("repo", "pkg", "name"),

	)

	a := NewPackageConfig(nil)
	if err := a.ParseDirectives("", withDirectives(
		"proto_plugin", "fake_proto label @fake//proto/plugin",
		"proto_rule", "fake_proto_library enabled true",
		"proto_language", "fake plugin fake_proto",
		"proto_language", "fake rule fake_proto_library",
	)); err != nil {
		t.Fatal(err)
	}
	b := a.Clone()

	initialState(t, a)
	initialState(t, b)

	b.rules["fake_proto_library"].Enabled = false

	initialState(t, a)
	finalState(t, b)
}

func testDirectives(t *testing.T, cases map[string]packageConfigTestCase) {
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			cfg := NewPackageConfig(nil)
			// t.Logf("test case: %+v", tc)
			if err := cfg.ParseDirectives(tc.rel, tc.directives); err != nil {
				if tc.err == nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tc.err.Error() != err.Error() {
					t.Fatalf("unexpected error: want %v, got %v", tc.err, err)
				}
				return
			}
			tc.check(t, cfg)
		})
	}
}

func allPackageChecks(checks ...packageConfigCheck) packageConfigCheck {
	return func(t *testing.T, cfg *PackageConfig) {
		for _, check := range checks {
			check(t, cfg)
		}
	}
}

func withDirectives(items ...string) (d []rule.Directive) {
	if len(items)%2 != 0 {
		panic("directive list must be a sequence of key/value pairs")
	}
	if len(items) < 2 {
		return
	}
	for i := 1; i < len(items); i = i + 2 {
		d = append(d, rule.Directive{Key: items[i-1], Value: items[i]})
	}
	return
}

func withImportpathPrefix(prefix string) packageConfigCheck {
	return func(t *testing.T, cfg *PackageConfig) {
		if cfg.importpathPrefix != prefix {
			t.Errorf("expected importpath prefix %s, got %s", prefix, cfg.importpathPrefix)
		}
	}
}
