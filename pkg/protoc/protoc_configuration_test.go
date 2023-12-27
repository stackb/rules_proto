package protoc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMergeSources(t *testing.T) {
	for name, tc := range map[string]struct {
		workDir           string
		rel               string
		plugins           []*PluginConfiguration
		stripImportPrefix string
		wantOutputs       []string
		wantMappings      map[string]string
	}{
		"empty": {
			wantOutputs:  []string{},
			wantMappings: map[string]string{},
		},
		"root package, simple case": {
			rel: "",
			plugins: withPluginConfigurations(&PluginConfiguration{
				Outputs: listOf("foo.py"),
			}),
			wantOutputs:  listOf("foo.py"),
			wantMappings: map[string]string{},
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
			wantOutputs:  listOf("foo.py"),
			wantMappings: map[string]string{},
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
		"child package, with strip_import_prefix": {
			rel: "test/proto",
			plugins: withPluginConfigurations(&PluginConfiguration{
				Outputs: listOf("test/proto/foo.py"),
			}),
			stripImportPrefix: "/test",
			wantOutputs:       listOf("foo.py"),
			wantMappings:      map[string]string{"foo.py": "/proto/foo.py"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			srcs, mappings := mergeSources(tc.workDir, tc.rel, tc.plugins, tc.stripImportPrefix)

			if diff := cmp.Diff(tc.wantOutputs, srcs); diff != "" {
				t.Errorf("srcs (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.wantMappings, mappings); diff != "" {
				t.Errorf("mappings (-want +got):\n%s", diff)
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
