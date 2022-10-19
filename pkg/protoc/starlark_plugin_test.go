package protoc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/google/go-cmp/cmp"
)

func TestLoadStarlarkPlugin(t *testing.T) {
	for name, tc := range map[string]struct {
		code        string
		ctx         *PluginContext
		wantErr     error
		wantPrinted string
		want        *PluginConfiguration
	}{
		"degenerate": {
			wantErr: fmt.Errorf(`test.star: plugin "test" was never declared`),
		},
		"wrong plugin name": {
			code: `
protoc.Plugin(
	name = "not-test",
	configure = lambda ctx: None,
)
			`,
			wantErr: fmt.Errorf(`test.star: plugin "test" was never declared`),
		},
		"missing configure attribute": {
			code: `
protoc.Plugin(
	name = "test", 
)
			`,
			wantErr: fmt.Errorf(`test.star: eval: Plugin: missing argument for configure`),
		},
		"configure attribute not callable": {
			code: `
protoc.Plugin(
	name = "test", 
	configure = "not-callable",
)
			`,
			wantErr: fmt.Errorf(`test.star: eval: Plugin: for parameter "configure": got string, want callable`),
		},
		"simple": {
			code: `
def configure(ctx):
	print(ctx)
	return protoc.PluginConfiguration(
		label = "//%s:python_plugin" % ctx.rel,
		outputs = ["foo.py", "bar.py"],
	)
    
protoc.Plugin(
	name = "test", 
	configure = configure,
)
`,
			ctx: &PluginContext{
				Rel: "mypkg",
			},
			want: &PluginConfiguration{
				Label:   label.New("", "mypkg", "python_plugin"),
				Outputs: []string{"foo.py", "bar.py"},
				Options: []string{},
			},
			wantPrinted: `PluginContext(package_config = PackageConfig(config = Config(repo_name = "", repo_root = "", work_dir = "")), plugin_config = LanguagePluginConfig(deps = [], enabled = False, implementation = "", label = "", name = "", options = []), proto_library = ProtoLibrary(base_name = "", deps = [], files = [], imports = [], name = "", srcs = [], strip_import_prefix = ""), rel = "mypkg")` + "\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			var err error
			var gotPrinted strings.Builder
			var plugin Plugin
			plugin, err = loadStarlarkPlugin("test", "test.star", strings.NewReader(tc.code), func(msg string) {
				gotPrinted.WriteString(msg)
				gotPrinted.Write([]byte{'\n'})
			}, func(configureErr error) {
				err = configureErr
			})
			if err != nil {
				if tc.wantErr != nil {
					if diff := cmp.Diff(tc.wantErr.Error(), err.Error()); diff != "" {
						t.Fatalf("StarlarkPlugin.Configure error (-want +got):\n%s", diff)
					}
					return
				} else {
					t.Fatalf("StarlarkPlugin.Configure error: %v", err)
				}
			}

			got := plugin.Configure(tc.ctx)
			t.Log(gotPrinted.String())
			if diff := cmp.Diff(tc.wantPrinted, gotPrinted.String()); diff != "" {
				t.Errorf("StarlarkPlugin.Configure print (-want +got):\n%s", diff)
			}

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("StarlarkPlugin.Configure (-want +got):\n%s", diff)
			}
		})
	}
}
