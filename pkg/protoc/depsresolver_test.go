package protoc

import (
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
)

func TestStripRel(t *testing.T) {
	for name, tc := range map[string]struct {
		rel  string
		in   string
		want string
	}{
		"empty": {
			rel:  "",
			in:   "",
			want: "",
		},
		"match": {
			rel:  "proto",
			in:   "proto/foo.txt",
			want: "foo.txt",
		},
		"no-match": {
			rel:  "proto",
			in:   "foo/proto.txt",
			want: "foo/proto.txt",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := StripRel(tc.rel, tc.in)
			if got != tc.want {
				t.Errorf("striprel: want %s, got %s", tc.want, got)
			}
		})
	}
}

func TestIsSameImport(t *testing.T) {
	for name, tc := range map[string]struct {
		from, to label.Label
		config   *config.Config
		want     bool
	}{
		"same": {
			from:   label.New("", "pkg", "a"),
			to:     label.New("", "pkg", "a"),
			config: &config.Config{RepoName: ""},
			want:   true,
		},
		"different": {
			from:   label.New("", "pkg", "a"),
			to:     label.New("", "pkg", "b"),
			config: &config.Config{RepoName: ""},
			want:   false,
		},
		"same repo": {
			from:   label.New("foo", "pkg", "a"),
			to:     label.New("", "pkg", "a"),
			config: &config.Config{RepoName: "foo"},
			want:   true,
		},
		"different repo": {
			from:   label.New("foo", "pkg", "a"),
			to:     label.New("bar", "pkg", "a"),
			config: &config.Config{RepoName: "foo"},
			want:   false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := isSameImport(tc.config, tc.from, tc.to)
			if got != tc.want {
				t.Errorf("isSameImport: want %t, got %t", tc.want, got)
			}
		})
	}
}
