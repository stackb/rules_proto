package protoc

import "testing"

func TestMergeSources(t *testing.T) {
	for name, tc := range map[string]struct {
		rel          string
		plugins      []*PluginConfiguration
		wantSrcs     []string
		wantMappings map[string]string
	}{
		"empty": {},
		"root package, simple case": {
			rel: "",
			plugins: withPluginConfigurations(&PluginConfiguration{
				Srcs: withSrcs("foo.py"),
			}),
			wantSrcs: withSrcs("foo.py"),
		},
		"root package, go_package case": {
			rel: "",
			plugins: withPluginConfigurations(&PluginConfiguration{
				Srcs:     withSrcs("foo.py"),
				Mappings: map[string]string{"foo.py": "com/github/example/foo.py"},
			}),
			wantSrcs:     withSrcs("foo.py"),
			wantMappings: map[string]string{"foo.py": "com/github/example/foo.py"},
		},
		"child package, simple case": {
			rel: "test/proto",
			plugins: withPluginConfigurations(&PluginConfiguration{
				Srcs: withSrcs("test/proto/foo.py"),
			}),
			wantSrcs: withSrcs("foo.py"),
		},
		"child package, mapped case": {
			rel: "test/proto",
			plugins: withPluginConfigurations(&PluginConfiguration{
				Srcs: withSrcs("foo.py"),
			}),
			wantSrcs:     withSrcs("foo.py"),
			wantMappings: map[string]string{"foo.py": "foo.py"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			srcs, mappings := mergeSources(tc.rel, tc.plugins)
			if len(tc.wantSrcs) != len(srcs) {
				t.Fatalf("srcs: want %d, got %d", len(tc.wantSrcs), len(srcs))
			}
			if len(tc.wantMappings) != len(mappings) {
				t.Fatalf("mappings: want %d, got %d", len(tc.wantMappings), len(mappings))
			}
			for i, got := range srcs {
				want := tc.wantSrcs[i]
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

func withSrcs(ss ...string) []string {
	return ss
}
