package protoc

import "testing"

func TestMergeSources(t *testing.T) {
	for name, tc := range map[string]struct {
		workDir      string
		rel          string
		plugins      []*PluginConfiguration
		wantOutputs  []string
		wantMappings map[string]string
	}{
		"empty": {},
		"root package, simple case": {
			rel: "",
			plugins: withPluginConfigurations(&PluginConfiguration{
				Outputs: listOf("foo.py"),
			}),
			wantOutputs: listOf("foo.py"),
		},
		"root package, go_package case": {
			rel: "",
			plugins: withPluginConfigurations(&PluginConfiguration{
				Outputs:  listOf("foo.py"),
				Mappings: map[string]string{"foo.py": "com/github/example/foo.py"},
			}),
			wantOutputs:  listOf("foo.py"),
			wantMappings: map[string]string{"foo.py": "com/github/example/foo.py"},
		},
		"child package, simple case": {
			rel: "test/proto",
			plugins: withPluginConfigurations(&PluginConfiguration{
				Outputs: listOf("test/proto/foo.py"),
			}),
			wantOutputs: listOf("foo.py"),
		},
		"child package, mapped case": {
			rel: "test/proto",
			plugins: withPluginConfigurations(&PluginConfiguration{
				Outputs: listOf("foo.py"),
			}),
			wantOutputs:  listOf("foo.py"),
			wantMappings: map[string]string{"foo.py": "foo.py"},
		},
		"external workspace, mapped case": {
			workDir: "/path/to/external/googleapis",
			rel:     "test/proto",
			plugins: withPluginConfigurations(&PluginConfiguration{
				Outputs: listOf("foo.py"),
			}),
			wantOutputs:  listOf("foo.py"),
			wantMappings: map[string]string{"foo.py": "foo.py"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			srcs, mappings := mergeSources(tc.workDir, tc.rel, tc.plugins)
			if len(tc.wantOutputs) != len(srcs) {
				t.Fatalf("srcs: want %d, got %d", len(tc.wantOutputs), len(srcs))
			}
			if len(tc.wantMappings) != len(mappings) {
				t.Fatalf("mappings: want %d, got %d", len(tc.wantMappings), len(mappings))
			}
			for i, got := range srcs {
				want := tc.wantOutputs[i]
				if want != got {
					t.Errorf("srcs %d: want %s, got %s", i, want, got)
				}
			}
			for name, got := range mappings {
				want := tc.wantMappings[name]
				if want != got {
					t.Errorf("mappings %q: want %s, got %s", name, want, got)
				}
			}
		})
	}
}

func withPluginConfigurations(cc ...*PluginConfiguration) []*PluginConfiguration {
	return cc
}

func listOf(ss ...string) []string {
	return ss
}
