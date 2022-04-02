// gencopy is a utility program that copies bazel outputs back into the
// workspace source tree.  Ideally, you don't have any generated files committed
// to VCS, but sometimes you do.
//
package main

import (
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParseFileMode(t *testing.T) {
	for name, tc := range map[string]struct {
		in   string
		want os.FileMode
	}{
		"ModePerm": {
			in:   "0777",
			want: os.FileMode(os.ModePerm),
		},
		"0644": {
			in:   "0644",
			want: os.FileMode(0o644),
		},
		"ModeSetgid": {
			in:   "020000000",
			want: os.FileMode(os.ModeSetgid),
		},
	} {
		t.Run(name, func(t *testing.T) {
			mode, err := parseFileMode(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			got := os.FileMode(mode)
			if tc.want != got {
				t.Errorf("want %#o, got %#o", tc.want, got)
			}
		})
	}
}

func TestMakePkgSrcDstPair(t *testing.T) {
	for name, tc := range map[string]struct {
		cfg      Config
		pkg      PackageConfig
		src, dst string
		want     SrcDst
	}{
		"degenerate": {},
		"WorkspaceRootDirectory": {
			cfg:  Config{WorkspaceRootDirectory: "/home"},
			src:  "file.txt",
			dst:  "file.txt",
			want: SrcDst{Src: "file.txt", Dst: "/home/file.txt"},
		},
		"TargetWorkspaceRoot": {
			cfg:  Config{WorkspaceRootDirectory: "/home"},
			pkg:  PackageConfig{TargetWorkspaceRoot: "external/foo"},
			src:  "../foo/file.txt",
			dst:  "file.txt",
			want: SrcDst{Src: "external/foo/file.txt", Dst: "/home/external/foo/file.txt"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := makePkgSrcDstPair(&tc.cfg, &tc.pkg, tc.src, tc.dst)

			if diff := cmp.Diff(tc.want, *got); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
