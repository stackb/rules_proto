// gencopy is a utility program that copies bazel outputs back into the
// workspace source tree.  Ideally, you don't have any generated files committed
// to VCS, but sometimes you do.
package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/testtools"
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
		"WorkspaceSubDirectory": {
			cfg:  Config{WorkspaceRootDirectory: "/home"},
			src:  "subdir/file.txt",
			dst:  "subdir/file.txt",
			want: SrcDst{Src: "subdir/file.txt", Dst: "/home/subdir/file.txt"},
		},
		"TargetWorkspaceRoot": {
			cfg:  Config{WorkspaceRootDirectory: "/home"},
			pkg:  PackageConfig{TargetWorkspaceRoot: "external/foo"},
			src:  "../foo/file.txt",
			dst:  "file.txt",
			want: SrcDst{Src: "../foo/file.txt", Dst: "/home/external/foo/file.txt"},
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

func TestRunPkg(t *testing.T) {
	for name, tc := range map[string]struct {
		cfg         Config
		files, want []testtools.FileSpec
	}{
		"degenerate": {},
		"simple": {
			cfg: Config{
				Extension:             "",
				FileMode:              "0644",
				Mode:                  "update",
				UpdateTargetLabelName: "api_v1_python_compiled_sources.update",
				PackageConfigs: []*PackageConfig{
					{
						GeneratedFiles:      []string{"example/assets/api.pb.go.gen"},
						SourceFiles:         []string{"example/assets/api.pb.go"},
						TargetLabel:         "@@//example/assets:api_go_compiled_sources",
						TargetPackage:       "example/assets",
						TargetWorkspaceRoot: "",
					},
				},
			},
			files: []testtools.FileSpec{
				{
					Path:    "example/assets/api.pb.go.gen",
					Content: "# generated file",
				},
			},
			want: []testtools.FileSpec{
				{
					Path:    "example/assets/api.pb.go",
					Content: "# generated file",
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			dir, cleanup := testtools.CreateFiles(t, tc.files)
			defer cleanup()

			if err := os.Chdir(dir); err != nil {
				t.Fatal(err)
			}

			t.Logf("BEFORE:")
			if err := listFiles(t, "."); err != nil {
				t.Logf("listfiles error: %v", err)
			}

			if err := run(&tc.cfg); err != nil {
				t.Fatal(err)
			}

			t.Logf("AFTER:")
			if err := listFiles(t, "."); err != nil {
				t.Logf("listfiles error: %v", err)
			}

			testtools.CheckFiles(t, dir, tc.want)
		})
	}
}

// listFiles - convenience debugging function to log the files under a given dir
func listFiles(t *testing.T, dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Logf("%v\n", err)
			return err
		}
		if info.Mode()&os.ModeSymlink > 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			t.Logf("%s -> %s", path, link)
			return nil
		}

		t.Log(path)
		return nil
	})
}
