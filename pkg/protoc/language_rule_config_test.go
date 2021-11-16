package protoc

import "testing"

type languageRuleConfigCheck func(t *testing.T, cfg *LanguageRuleConfig)

func TestLanguageRuleConfigClone(t *testing.T) {
	a := NewLanguageRuleConfig(nil, "proto_compile")
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
		"proto_rule resolve": {
			directives: withDirectives(
				"proto_rule", "fake_proto_library resolve google/protobuf/([a-z]+).proto @org_golang_google_protobuf//types/known/$1pb",
			),
			check: withLanguageRule("fake_proto_library", withRuleResolvesEquals(
				"google/protobuf/([a-z]+).proto @org_golang_google_protobuf//types/known/$1pb",
			)),
		},
		"proto_rule attr 1": {
			directives: withDirectives(
				"proto_rule", "fake_proto_library attr count 1",
			),
			check: withLanguageRule("fake_proto_library", withRuleAttrEquals(
				"count",
				"1",
			)),
		},
		"proto_rule attr 2": {
			directives: withDirectives(
				"proto_rule", "fake_proto_library attr count 1",
				"proto_rule", "fake_proto_library attr count 2",
			),
			check: withLanguageRule("fake_proto_library", withRuleAttrEquals(
				"count",
				"1",
				"2",
			)),
		},
		"proto_rule attr -count": {
			directives: withDirectives(
				"proto_rule", "fake_proto_library attr count 1",
				"proto_rule", "fake_proto_library attr count 2",
				"proto_rule", "fake_proto_library attr -count 1",
			),
			check: withLanguageRule("fake_proto_library", withRuleAttrEquals(
				"count",
				"2",
			)),
		},
		"proto_rule -attr count": {
			directives: withDirectives(
				"proto_rule", "fake_proto_library attr count 1",
				"proto_rule", "fake_proto_library attr count 2",
				"proto_rule", "fake_proto_library -attr count",
			),
			check: withLanguageRule("fake_proto_library", withRuleAttrEquals(
				"count",
			)),
		},
		"proto_rule attr with space": {
			directives: withDirectives(
				"proto_rule", "fake_proto_library attr args --lib ES2015",
			),
			check: withLanguageRule("fake_proto_library", withRuleAttrEquals(
				"--lib ES2015",
			)),
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

func withRuleResolvesEquals(resolves ...string) languageRuleConfigCheck {
	return func(t *testing.T, cfg *LanguageRuleConfig) {
		got := cfg.GetRewrites()
		if len(resolves) != len(got) {
			t.Fatalf("rule resolves: want %d, got %d", len(resolves), len(got))
		}
		for i := 0; i < len(got); i++ {
			expected := resolves[i]
			actual := got[i]
			original := actual.Match.String() + " " + actual.Replace
			if expected != original {
				t.Errorf("rule resolve #%d: want %s, got %s", i, expected, original)
			}
		}
	}
}

func withRuleAttrEquals(name string, want ...string) languageRuleConfigCheck {
	return func(t *testing.T, cfg *LanguageRuleConfig) {
		got := cfg.GetAttr(name)
		if len(want) != len(got) {
			t.Fatalf("rule attr %s: want %d, got %d", name, len(want), len(got))
		}
		for i := 0; i < len(got); i++ {
			expected := want[i]
			actual := got[i]
			if expected != actual {
				t.Errorf("rule attr %s #%d: want %s, got %s", name, i, expected, actual)
			}
		}
	}
}
