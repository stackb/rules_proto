package protobuf

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/google/go-cmp/cmp"
)

func TestGenerateRules(t *testing.T) {
	for name, tc := range map[string]struct {
		args language.GenerateArgs
		want language.GenerateResult
	}{
		"empty case": {
			args: language.GenerateArgs{
				Config: makeTestConfig(""),
			},
			want: language.GenerateResult{
				Gen:     []*rule.Rule{},
				Empty:   []*rule.Rule{},
				Imports: []interface{}{},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			tmpdir := os.Getenv("TEST_TMPDIR")
			dir, err := ioutil.TempDir(tmpdir, "")
			if err != nil {
				t.Fatalf("ioutil.TempDir(%q, %q) failed with %v; want success", tmpdir, "", err)
			}
			defer os.RemoveAll(dir)

			ext := NewProtobufLang("test")
			tc.args.Rel = dir

			got := ext.GenerateRules(tc.args)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Error("unexpected diff:", diff)
			}
		})
	}
}

func makeTestConfig(repoName string) *config.Config {
	return &config.Config{
		RepoName: repoName,
		Exts:     make(map[string]interface{}),
	}
}
