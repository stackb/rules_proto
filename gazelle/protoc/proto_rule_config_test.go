package protoc

import (
	"testing"

	"github.com/bazelbuild/bazel-gazelle/rule"
)

type ruleConfigCheck func(t *testing.T, cfg *ProtoRuleConfig)

func init() {
	MustRegisterProtoRule("fake_proto_library", &fakeProtoLibrary{})
}

func TestProtoRuleConfigClone(t *testing.T) {
	a := newProtoRuleConfig("proto_compile")
	a.Deps = map[string]bool{
		"d1": true,
		"d2": true,
	}
	b := a.Clone()

	check := allRuleChecks(
		withProtoRuleName("proto_compile"),
		withProtoRuleEnabled(true),
		withRuleDepsEquals("d1", "d2"),
	)

	check(t, a)
	check(t, b)

	b.Enabled = false

	withProtoRuleEnabled(true)(t, a)
	withProtoRuleEnabled(false)(t, b)
}

func TestProtoRuleDirectives(t *testing.T) {
	testDirectives(t, map[string]packageConfigTestCase{
		"proto_rule enabled": {
			directives: withDirectives("proto_rule", "fake_proto_library enabled true"),
			check:      withProtoRule("fake_proto_library", withProtoRuleEnabled(true)),
		},
		"proto_rule dep": {
			directives: withDirectives(
				"proto_rule", "fake_proto_library dep @fake//lib",
			),
			check: withProtoRule("fake_proto_library", withRuleDepsEquals("@fake//lib")),
		},
		"proto_rule -dep": {
			directives: withDirectives(
				"proto_rule", "fake_proto_library +dep @fake//lib",
				"proto_rule", "fake_proto_library -dep @fake//lib",
			),
			check: withProtoRule("fake_proto_library", withRuleDepsEquals()),
		},
	})
}

func allRuleChecks(checks ...ruleConfigCheck) ruleConfigCheck {
	return func(t *testing.T, cfg *ProtoRuleConfig) {
		for _, check := range checks {
			check(t, cfg)
		}
	}
}

func withProtoRule(name string, checks ...ruleConfigCheck) packageConfigCheck {
	return func(t *testing.T, cfg *ProtoPackageConfig) {
		rule, ok := cfg.rules[name]
		if !ok {
			t.Fatal("rule not found", name)
		}
		for _, check := range checks {
			check(t, rule)
		}
	}
}

func withProtoRuleEnabled(enabled bool) ruleConfigCheck {
	return func(t *testing.T, cfg *ProtoRuleConfig) {
		if cfg.Enabled != enabled {
			t.Errorf("rule label: want %t, got %t", enabled, cfg.Enabled)
		}
	}
}

func withProtoRuleName(name string) ruleConfigCheck {
	return func(t *testing.T, cfg *ProtoRuleConfig) {
		if cfg.Name != name {
			t.Errorf("rule name: want %s, got %s", name, cfg.Name)
		}
	}
}

func withRuleDepsEquals(deps ...string) ruleConfigCheck {
	return func(t *testing.T, cfg *ProtoRuleConfig) {
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

// fakeProtoLibrary implements a mock ProtoRule
type fakeProtoLibrary struct {
}

// KindInfo implements part of the ProtoRule interface.
func (s *fakeProtoLibrary) KindInfo() rule.KindInfo {
	return rule.KindInfo{}
}

// LoadInfo implements part of the ProtoRule interface.
func (s *fakeProtoLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{}
}

// GenerateRule implements part of the ProtoRule interface.
func (s *fakeProtoLibrary) GenerateRule(rc *ProtoRuleConfig, pc *ProtocConfiguration) RuleProvider {
	return nil
}
