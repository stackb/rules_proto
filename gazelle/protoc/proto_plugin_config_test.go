package protoc

import (
	"testing"

	"github.com/bazelbuild/bazel-gazelle/label"
)

type pluginConfigCheck func(t *testing.T, cfg *ProtoPluginConfig)

func TestProtoPluginDirectives(t *testing.T) {
	testDirectives(t, map[string]packageConfigTestCase{
		"proto_plugin label": {
			directives: withDirectives("proto_plugin", "py_proto label //protoc:py_proto_plugin"),
			check:      withProtoPlugin("py_proto", withPluginLabelEquals("", "protoc", "py_proto_plugin")),
		},
		"proto_plugin tool": {
			directives: withDirectives("proto_plugin", "gogofast tool @com_github_gogo_protobuf//protoc-gen-gogofast"),
			check:      withProtoPlugin("gogofast", withPluginToolEquals("com_github_gogo_protobuf", "protoc-gen-gogofast", "protoc-gen-gogofast")),
		},
		"proto_plugin option": {
			directives: withDirectives("proto_plugin", "gogofast option grpc"),
			check:      withProtoPlugin("gogofast", withPluginOptionsEquals("grpc")),
		},
		"proto_plugin +option": {
			directives: withDirectives("proto_plugin", "gogofast +option grpc"),
			check:      withProtoPlugin("gogofast", withPluginOptionsEquals("grpc")),
		},
		"proto_plugin -option": {
			directives: withDirectives(
				"proto_plugin", "gogofast +option grpc",
				"proto_plugin", "gogofast -option grpc",
			),
			check: withProtoPlugin("gogofast", withPluginOptionsEquals()),
		},
	})
}

func withProtoPlugin(name string, checks ...pluginConfigCheck) packageConfigCheck {
	return func(t *testing.T, cfg *ProtoPackageConfig) {
		plugin, ok := cfg.plugins[name]
		if !ok {
			t.Fatal("plugin not found", name)
		}
		for _, check := range checks {
			check(t, plugin)
		}
	}
}

func withPluginLabelEquals(repo, pkg, name string) pluginConfigCheck {
	return func(t *testing.T, cfg *ProtoPluginConfig) {
		want := label.Label{Repo: repo, Pkg: pkg, Name: name}
		got := cfg.Label
		if want.String() != got.String() {
			t.Errorf("plugin label: want %s, got %s", want, got)
		}
	}
}

func withPluginToolEquals(repo, pkg, name string) pluginConfigCheck {
	return func(t *testing.T, cfg *ProtoPluginConfig) {
		want := label.Label{Repo: repo, Pkg: pkg, Name: name}
		got := cfg.Tool
		if want.String() != got.String() {
			t.Errorf("plugin label: want %s, got %s", want, got)
		}
	}
}

func withPluginOptionsEquals(opts ...string) pluginConfigCheck {
	return func(t *testing.T, cfg *ProtoPluginConfig) {
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
