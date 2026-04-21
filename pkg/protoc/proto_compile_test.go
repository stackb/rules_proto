package protoc

import (
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/emicklei/proto"
)

func TestHasOverlap(t *testing.T) {
	for name, tc := range map[string]struct {
		a, b []string
		want bool
	}{
		"both empty": {
			a: nil, b: nil, want: false,
		},
		"no overlap": {
			a: []string{"a", "b"}, b: []string{"c", "d"}, want: false,
		},
		"overlap": {
			a: []string{"a", "b"}, b: []string{"b", "c"}, want: true,
		},
		"identical": {
			a: []string{"a"}, b: []string{"a"}, want: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			if got := hasOverlap(tc.a, tc.b); got != tc.want {
				t.Errorf("hasOverlap(%v, %v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

// packageLevelPlugin is a plugin that produces a single output per proto
// package, regardless of which proto file it's configured with. This simulates
// plugins like protoc-gen-prost.
type packageLevelPlugin struct{}

func (p *packageLevelPlugin) Name() string {
	return "protoc:package_level"
}

func (p *packageLevelPlugin) Configure(ctx *PluginContext) *PluginConfiguration {
	return &PluginConfiguration{
		Label:   ctx.PluginConfig.Label,
		Outputs: []string{"my_package.rs"},
	}
}

func init() {
	Plugins().MustRegisterPlugin(&packageLevelPlugin{})
}

func makeProtoLibrary(name, filename string) ProtoLibrary {
	r := rule.NewRule("proto_library", name)
	f := NewFile("pkg", filename)
	f.pkg = proto.Package{Name: "my.package"}
	f.messages = append(f.messages, proto.Message{Name: "Msg"})
	return NewOtherProtoLibrary(nil, r, f)
}

func aggregationPackageConfig() *PackageConfig {
	c := NewPackageConfig(&config.Config{})
	if err := c.ParseDirectives("pkg", withDirectives(
		"proto_rule", "proto_compile implementation stackb:rules_proto:proto_compile",
		"proto_plugin", "pkg_plugin implementation protoc:package_level",
		"proto_plugin", "pkg_plugin enabled true",
		"proto_language", "rust plugin pkg_plugin",
		"proto_language", "rust enabled true",
		"proto_language", "rust rule proto_compile",
	)); err != nil {
		panic("bad config: " + err.Error())
	}
	return c
}

// ExamplePackageAggregation demonstrates that when two proto_library rules
// produce the same output file, their proto_compile rules are merged into a
// single rule using the "protos" attribute.
func ExamplePackage_aggregation() {
	pkg := NewPackage(
		"pkg",
		aggregationPackageConfig(),
		makeProtoLibrary("a_proto", "a.proto"),
		makeProtoLibrary("b_proto", "b.proto"),
	)
	formaatRules(pkg.Rules()...)
	// Output:
	// proto_compile(
	//     name = "a_rust_compile",
	//     output_mappings = ["my_package.rs=my_package.rs"],
	//     outputs = ["my_package.rs"],
	//     plugins = ["//:"],
	//     protos = [
	//         "a_proto",
	//         "b_proto",
	//     ],
	// )
}

// ExamplePackageNoAggregation demonstrates that when two proto_library rules
// produce different output files, separate proto_compile rules are emitted.
func ExamplePackage_noAggregation() {
	pkg := NewPackage(
		exampleDir,
		examplePackageConfig(),
		exampleProtoLibrary(),
	)
	formaatRules(pkg.Rules()...)
	// Output:
	// proto_compile(
	//     name = "test_fake_compile",
	//     outputs = ["test_fake.pb.go"],
	//     plugins = ["@build_stack_rules_proto//plugin/builtin:fake"],
	//     proto = "test_proto",
	// )
}
