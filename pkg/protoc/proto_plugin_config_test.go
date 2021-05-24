package protoc

import (
	"testing"

	"github.com/bazelbuild/bazel-gazelle/label"
)

type LanguagePluginConfigCheck func(t *testing.T, cfg *LanguagePluginConfig)

func TestPluginDirectives(t *testing.T) {
	testDirectives(t, map[string]packageConfigTestCase{
		"proto_plugin label": {
			directives: withDirectives("proto_plugin", "fake_proto label //protoc:fake_proto_plugin"),
			check:      withPlugin("fake_proto", withPluginLabelEquals("", "protoc", "fake_proto_plugin")),
		},
		"proto_plugin option": {
			directives: withDirectives("proto_plugin", "fake_proto option grpc"),
			check:      withPlugin("fake_proto", withPluginOptionsEquals("grpc")),
		},
		"proto_plugin +option": {
			directives: withDirectives("proto_plugin", "fake_proto +option grpc"),
			check:      withPlugin("fake_proto", withPluginOptionsEquals("grpc")),
		},
		"proto_plugin -option": {
			directives: withDirectives(
				"proto_plugin", "fake_proto +option grpc",
				"proto_plugin", "fake_proto -option grpc",
			),
			check: withPlugin("fake_proto", withPluginOptionsEquals()),
		},
	})
}

func withPlugin(name string, checks ...LanguagePluginConfigCheck) packageConfigCheck {
	return func(t *testing.T, cfg *PackageConfig) {
		plugin, ok := cfg.plugins[name]
		if !ok {
			t.Fatal("plugin not found", name)
		}
		for _, check := range checks {
			check(t, plugin)
		}
	}
}

func withPluginLabelEquals(repo, pkg, name string) LanguagePluginConfigCheck {
	return func(t *testing.T, cfg *LanguagePluginConfig) {
		want := label.New(repo, pkg, name)
		got := cfg.Label
		if want.String() != got.String() {
			t.Errorf("plugin label: want %s, got %s", want, got)
		}
	}
}

func withPluginOptionsEquals(opts ...string) LanguagePluginConfigCheck {
	return func(t *testing.T, cfg *LanguagePluginConfig) {
		got := cfg.GetOptions()
		if len(opts) != len(got) {
			t.Fatalf("plugin options: want %d, got %d", len(opts), len(got))
		}
		for i := 0; i < len(got); i++ {
			expected := opts[i]
			actual := got[i]
			if expected != actual {
				t.Errorf("plugin option #%d: want %s, got %s", i, expected, actual)
			}
		}
	}
}
