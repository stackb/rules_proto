package rules_python

import (
	"testing"

	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/google/go-cmp/cmp"
)

func TestMaybeStripImportPrefix(t *testing.T) {
	for name, tc := range map[string]struct {
		specs             []resolve.ImportSpec
		stripImportPrefix string
		want              []resolve.ImportSpec
	}{
		"degenerate": {},
		"empty string": {
			specs: []resolve.ImportSpec{
				{Imp: "foo/bar/baz.proto"},
			},
			want: []resolve.ImportSpec{
				{Imp: "foo/bar/baz.proto"},
			},
		},
		"non-matching-prefix": {
			specs: []resolve.ImportSpec{
				{Imp: "foo/bar/baz.proto"},
			},
			stripImportPrefix: "nope",
			want: []resolve.ImportSpec{
				{Imp: "foo/bar/baz.proto"},
			},
		},
		"matching-prefix-absolute": {
			specs: []resolve.ImportSpec{
				{Imp: "foo/bar/baz.proto"},
			},
			stripImportPrefix: "/foo/bar",
			want: []resolve.ImportSpec{
				{Imp: "baz.proto"},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := maybeStripImportPrefix(tc.specs, tc.stripImportPrefix)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
