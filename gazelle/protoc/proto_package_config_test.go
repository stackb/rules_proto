package protoc

import (
	"testing"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

type packageConfigCheck func(t *testing.T, cfg *ProtoPackageConfig)

type packageConfigTestCase struct {
	rel        string
	directives []rule.Directive
	check      packageConfigCheck
	err        error
}

func TestProtoPackageConfigClone(t *testing.T) {
	initialState := allPackageChecks(
		withProtoLang("py",
			withProtoLangEnabled(true),
		),
		withProtoRule("fake_proto_library",
			withProtoRuleEnabled(true),
		),
	)
	finalState := allPackageChecks(
		withProtoLang("py",
			withProtoLangEnabled(true),
		),
		withProtoRule("fake_proto_library",
			withProtoRuleEnabled(false),
		),
		withProtoPlugin("py_proto",
			withPluginToolEquals("other", "some", "tool"),
		),
	)

	a := newProtoPackageConfig()
	a.parseDirectives("", withDirectives(
		"proto_plugin", "py_proto label @fake//proto/plugin",
		"proto_rule", "fake_proto_library enabled true",
		"proto_lang", "py label @py//proto/lang",
		"proto_lang", "py plugin py_proto",
		"proto_lang", "py rule fake_proto_library",
	))
	b := a.Clone()

	initialState(t, a)
	initialState(t, b)

	b.rules["fake_proto_library"].Enabled = false
	b.plugins["py_proto"].Tool = label.Label{Repo: "other", Pkg: "some", Name: "tool"}

	initialState(t, a)
	finalState(t, b)
}

func testDirectives(t *testing.T, cases map[string]packageConfigTestCase) {
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			cfg := newProtoPackageConfig()
			t.Logf("test case: %+v", tc)
			if err := cfg.parseDirectives(tc.rel, tc.directives); err != nil {
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
	return func(t *testing.T, cfg *ProtoPackageConfig) {
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
	return func(t *testing.T, cfg *ProtoPackageConfig) {
		if cfg.importpathPrefix != prefix {
			t.Errorf("expected importpath prefix %s, got %s", prefix, cfg.importpathPrefix)
		}
	}
}
