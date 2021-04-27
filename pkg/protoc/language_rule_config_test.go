package protoc

import (
	"testing"
)

type languageRuleConfigCheck func(t *testing.T, cfg *LanguageRuleConfig)

func TestLanguageRuleConfigClone(t *testing.T) {
	a := newLanguageRuleConfig("proto_compile")
	a.Deps = map[string]bool{
		"d1": true,
		"d2": true,
	}
	b := a.clone()

	check := allRuleChecks(
		withLanguageRuleName("proto_compile"),
		withLanguageRuleEnabled(true),
		withRuleDepsEquals("d1", "d2"),
	)

	check(t, a)
	check(t, b)

	b.Enabled = false

	withLanguageRuleEnabled(true)(t, a)
	withLanguageRuleEnabled(false)(t, b)
}

func TestRuleDirectives(t *testing.T) {
	testDirectives(t, map[string]packageConfigTestCase{
		"proto_rule enabled": {
			directives: withDirectives("proto_rule", "fake_proto_library enabled true"),
			check:      withLanguageRule("fake_proto_library", withLanguageRuleEnabled(true)),
		},
		"proto_rule dep": {
			directives: withDirectives(
				"proto_rule", "fake_proto_library dep @fake//lib",
			),
			check: withLanguageRule("fake_proto_library", withRuleDepsEquals("@fake//lib")),
		},
		"proto_rule -dep": {
			directives: withDirectives(
				"proto_rule", "fake_proto_library +dep @fake//lib",
				"proto_rule", "fake_proto_library -dep @fake//lib",
			),
			check: withLanguageRule("fake_proto_library", withRuleDepsEquals()),
		},
	})
}

func allRuleChecks(checks ...languageRuleConfigCheck) languageRuleConfigCheck {
	return func(t *testing.T, cfg *LanguageRuleConfig) {
		for _, check := range checks {
			check(t, cfg)
		}
	}
}

func withLanguageRule(name string, checks ...languageRuleConfigCheck) packageConfigCheck {
	return func(t *testing.T, cfg *PackageConfig) {
		rule, ok := cfg.rules[name]
		if !ok {
			t.Fatal("rule not found", name)
		}
		for _, check := range checks {
			check(t, rule)
		}
	}
}

func withLanguageRuleEnabled(enabled bool) languageRuleConfigCheck {
	return func(t *testing.T, cfg *LanguageRuleConfig) {
		if cfg.Enabled != enabled {
			t.Errorf("rule label: want %t, got %t", enabled, cfg.Enabled)
		}
	}
}

func withLanguageRuleName(name string) languageRuleConfigCheck {
	return func(t *testing.T, cfg *LanguageRuleConfig) {
		if cfg.Name != name {
			t.Errorf("rule name: want %s, got %s", name, cfg.Name)
		}
	}
}

func withRuleDepsEquals(deps ...string) languageRuleConfigCheck {
	return func(t *testing.T, cfg *LanguageRuleConfig) {
		got := cfg.GetDeps()
		if len(deps) != len(got) {
			t.Fatalf("rule deps: want %d, got %d", len(deps), len(got))
		}
		for i := 0; i < len(got); i++ {
			expected := deps[i]
			actual := got[i]
			if expected != actual {
				t.Errorf("rule dep #%d: want %s, got %s", i, expected, actual)
			}
		}
	}
}
