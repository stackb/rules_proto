package protoc

import (
	"testing"

	"github.com/bazelbuild/bazel-gazelle/rule"
)

type configCheck func(t *testing.T, cfg *protoPackageConfig)

func TestClone(t *testing.T) {
	cfg := newProtoPackageConfig()
	cfg.rules["foo"] = &protoRuleConfig{"foo", true}
	cfg.importpathPrefix = "github.com/foo/bar"

	check := all(
		hasLanguageConfig("py", true),
		hasRuleExclusion("foo", true),
		hasImportpathPrefix("github.com/foo/bar"),
	)

	check(t, cfg)
	check(t, cfg.Clone())
}

func TestParseDirectives(t *testing.T) {
	type testCase struct {
		rel        string
		directives []rule.Directive
		check      configCheck
	}

	for name, tc := range map[string]testCase{
		"proto_rule not mentioned": {
			directives: []rule.Directive{},
			check: all(
				hasRuleInclusion("foo", false),
				hasRuleExclusion("foo", false),
			),
		},
		"proto_rule present": {
			directives: []rule.Directive{
				{
					Key:   "proto_rule",
					Value: "foo",
				},
			},
			check: all(
				hasRuleInclusion("foo", true),
				hasRuleExclusion("foo", false),
			),
		},
		"proto_rule positive": {
			directives: []rule.Directive{
				{
					Key:   "proto_rule",
					Value: "+foo",
				},
			},
			check: all(
				hasRuleInclusion("foo", true),
				hasRuleExclusion("foo", false),
			),
		},
		"proto_rule negative": {
			directives: []rule.Directive{
				{
					Key:   "proto_rule",
					Value: "-foo",
				},
			},
			check: all(
				hasRuleInclusion("foo", false),
				hasRuleExclusion("foo", true),
			),
		},
		"proto_rule negative glob": {
			directives: []rule.Directive{
				{
					Key:   "proto_rule",
					Value: "-f*",
				},
			},
			check: all(
				hasRuleInclusion("foo", false),
				hasRuleExclusion("foo", true),
				hasRuleExclusion("baz", false),
				hasRuleExclusion("fdr", true),
			),
		},
		"proto_rule negative glob recover (glob before)": {
			directives: []rule.Directive{
				{
					Key:   "proto_rule",
					Value: "-f*",
				},
				{
					Key:   "proto_rule",
					Value: "+fdr",
				},
			},
			check: all(
				hasRuleInclusion("foo", false),
				hasRuleExclusion("foo", true),
				hasRuleExclusion("baz", false),
				hasRuleExclusion("fdr", false),
				hasRuleInclusion("fdr", true),
			),
		},
		"proto_rule negative glob recover (glob after)": {
			directives: []rule.Directive{
				{
					Key:   "proto_rule",
					Value: "+fdr",
				},
				{
					Key:   "proto_rule",
					Value: "-f*",
				},
			},
			check: all(
				hasRuleInclusion("foo", false),
				hasRuleExclusion("foo", true),
				hasRuleExclusion("baz", false),
				hasRuleExclusion("fdr", false),
				hasRuleInclusion("fdr", true),
			),
		},
		"proto_language": {
			directives: []rule.Directive{
				{
					Key:   "proto_language",
					Value: "py",
				},
			},
			check: hasLanguageConfig("py", true),
		},
		"+proto_language": {
			directives: []rule.Directive{
				{
					Key:   "proto_language",
					Value: "+py",
				},
			},
			check: hasLanguageConfig("py", true),
		},
		"-proto_language": {
			directives: []rule.Directive{
				{
					Key:   "proto_language",
					Value: "-py",
				},
			},
			check: hasLanguageConfig("py", false),
		},
		"-proto_language recover": {
			directives: []rule.Directive{
				{
					Key:   "proto_language",
					Value: "-py",
				},
				{
					Key:   "proto_language",
					Value: "+py",
				},
			},
			check: hasLanguageConfig("py", true),
		},
		"prefix": {
			directives: []rule.Directive{
				{
					Key:   "prefix",
					Value: "github.com/foo/bar",
				},
			},
			check: hasImportpathPrefix("github.com/foo/bar"),
		},
	} {
		t.Run(name, func(t *testing.T) {
			cfg := newProtoPackageConfig()
			cfg.parseDirectives(tc.rel, tc.directives)
			tc.check(t, cfg)
		})
	}
}

func all(checks ...configCheck) configCheck {
	return func(t *testing.T, cfg *protoPackageConfig) {
		for _, c := range checks {
			c(t, cfg)
		}
	}
}

func hasRuleExclusion(name string, expected bool) configCheck {
	return func(t *testing.T, cfg *protoPackageConfig) {
		actual := cfg.IsRuleExcluded(name)
		if expected != actual {
			t.Errorf("rule %q exclusion: expected %t", name, expected)
		}
	}
}

func hasRuleInclusion(name string, expected bool) configCheck {
	return func(t *testing.T, cfg *protoPackageConfig) {
		actual := cfg.IsRuleIncluded(name)
		if expected != actual {
			t.Errorf("rule %q inclusion: expected %t", name, expected)
		}
	}
}

func hasLanguageConfig(name string, present bool) configCheck {
	return func(t *testing.T, cfg *protoPackageConfig) {
		_, ok := cfg.languages[name]
		if ok && !present {
			t.Errorf("expected language to be excluded %s", name)
		}
		if !ok && present {
			t.Errorf("expected language to be included %s", name)
		}
	}
}

func hasImportpathPrefix(prefix string) configCheck {
	return func(t *testing.T, cfg *protoPackageConfig) {
		if cfg.importpathPrefix != prefix {
			t.Errorf("expected importpath prefix %s, got %s", prefix, cfg.importpathPrefix)
		}
	}
}
