package protoc

import (
	"testing"
)

type langConfigCheck func(t *testing.T, cfg *ProtoLangConfig)

func TestProtoLangConfigClone(t *testing.T) {
	a := newProtoLangConfig("py")
	a.Plugins["py_proto"] = newProtoPluginConfig("py_proto")
	a.Rules["py_proto_library"] = newProtoRuleConfig("py_proto_library")
	b := a.Clone()

	check := allLangChecks(
		withProtoLangEnabled(true),
	)

	check(t, a)
	check(t, b)

	b.Enabled = false
	b.Plugins["py_proto"].Enabled = false
	b.Rules["py_proto_library"].Enabled = false

	withProtoLangEnabled(true)(t, a)
	withProtoLangEnabled(false)(t, b)

	withProtoLangPluginEnabled("py_proto", true)(t, a)
	withProtoLangPluginEnabled("py_proto", false)(t, b)

	withProtoLangRuleEnabled("py_proto_library", true)(t, a)
	withProtoLangRuleEnabled("py_proto_library", false)(t, b)
}

func TestprotoLangDirectives(t *testing.T) {
	testDirectives(t, map[string]packageConfigTestCase{
		"proto_lang enabled": {
			directives: withDirectives("proto_lang", "foo enabled true"),
			check:      withProtoLang("foo", withProtoLangEnabled(true)),
		},
		"proto_lang disabled": {
			directives: withDirectives("proto_lang", "foo enabled false"),
			check:      withProtoLang("foo", withProtoLangEnabled(false)),
		},
		"proto_lang plugin": {
			directives: withDirectives(
				"proto_plugin", "py_proto label @foo//plugin",
				"proto_lang", "foo plugin py_proto",
			),
			check: withProtoLang("foo",
				withProtoLangEnabled(true),
				withProtoLangPluginEnabled("py_proto", true),
			),
		},
		"proto_lang -plugin": {
			directives: withDirectives(
				"proto_plugin", "py_proto label @foo//plugin",
				"proto_lang", "foo +plugin py_proto",
				"proto_lang", "foo -plugin py_proto",
			),
			check: withProtoLang("foo",
				withProtoLangEnabled(true),
				withProtoLangPluginEnabled("py_proto", false),
			),
		},
		"proto_lang rule": {
			directives: withDirectives(
				"proto_rule", "fake_proto_library enabled true",
				"proto_lang", "foo rule fake_proto_library",
			),
			check: withProtoLang("foo",
				withProtoLangEnabled(true),
				withProtoLangRuleEnabled("fake_proto_library", true),
			),
		},
		"proto_lang -rule": {
			directives: withDirectives(
				"proto_rule", "fake_proto_library enabled true",
				"proto_lang", "foo +rule fake_proto_library",
				"proto_lang", "foo -rule fake_proto_library",
			),
			check: withProtoLang("foo",
				withProtoLangEnabled(true),
				withProtoLangRuleEnabled("fake_proto_library", false),
			),
		},
	})
}

func allLangChecks(checks ...langConfigCheck) langConfigCheck {
	return func(t *testing.T, cfg *ProtoLangConfig) {
		for _, check := range checks {
			check(t, cfg)
		}
	}
}

func withProtoLang(name string, checks ...langConfigCheck) packageConfigCheck {
	return func(t *testing.T, cfg *ProtoPackageConfig) {
		lang, ok := cfg.langs[name]
		if !ok {
			t.Fatal("lang not found", name)
		}
		for _, check := range checks {
			check(t, lang)
		}
	}
}

func withProtoLangEnabled(enabled bool) langConfigCheck {
	return func(t *testing.T, cfg *ProtoLangConfig) {
		if cfg.Enabled != enabled {
			t.Logf("withProtoLangEnabled cfg: %+v", cfg)
			t.Errorf("lang enabled: want %t, got %t", enabled, cfg.Enabled)
		}
	}
}

func withProtoLangPluginEnabled(name string, enabled bool) langConfigCheck {
	return func(t *testing.T, cfg *ProtoLangConfig) {
		plugin, ok := cfg.Plugins[name]
		if !ok {
			t.Fatal("plugin not found:", name)
		}
		if plugin.Enabled != enabled {
			t.Logf("failing lang config: %+v", cfg)
			t.Errorf("lang plugin enabled: want %t, got %t", enabled, plugin.Enabled)
		}
	}
}

func withProtoLangRuleEnabled(name string, enabled bool) langConfigCheck {
	return func(t *testing.T, cfg *ProtoLangConfig) {
		rule, ok := cfg.Rules[name]
		if !ok {
			t.Fatal("rule not found:", name)
		}
		if rule.Enabled != enabled {
			t.Errorf("lang rule enabled: want %t, got %t", enabled, rule.Enabled)
		}
	}
}
