package starlarkrepository

import (
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/buildtools/build"
)

// readFileLines tests

func TestReadLines(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name: "filters empty lines and comments",
			input: `.build/
# This is a comment
.swiftpm/
.vscode/

`,
			expected: []string{".build/", ".swiftpm/", ".vscode/"},
		},
		{
			name: "handles only comments",
			input: `# Comment 1
# Comment 2
# Comment 3
`,
			expected: []string{},
		},
		{
			name: "handles mixed content",
			input: `  .build/
# Comment

.swiftpm/
   # Another comment
.vscode/
`,
			expected: []string{".build/", ".swiftpm/", ".vscode/"},
		},
		{
			name:     "empty file",
			input:    "",
			expected: []string{},
		},
		{
			name: "only whitespace",
			input: `


`,
			expected: []string{},
		},
		{
			name: "trims whitespace",
			input: `   .build/
	.swiftpm/
  .vscode/
`,
			expected: []string{".build/", ".swiftpm/", ".vscode/"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()
			testFile := filepath.Join(tmpDir, "test.txt")

			// Write input content
			if err := os.WriteFile(testFile, []byte(tt.input), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			// Create a no-op log function for tests
			noop := func(format string, args ...any) {}

			// Run the function
			lines, err := readFileLines(testFile, noop)
			if err != nil {
				t.Fatalf("readLines failed: %v", err)
			}

			// Handle nil vs empty slice
			if tt.expected == nil || len(tt.expected) == 0 {
				if len(lines) != 0 {
					t.Errorf("Expected empty result, got %d lines: %v", len(lines), lines)
				}
				return
			}

			// Compare results
			if len(lines) != len(tt.expected) {
				t.Errorf("Length mismatch: got %d lines, want %d lines", len(lines), len(tt.expected))
				t.Errorf("got: %v", lines)
				t.Errorf("want: %v", tt.expected)
				return
			}

			for i, line := range lines {
				if line != tt.expected[i] {
					t.Errorf("Line %d mismatch:\ngot:  %q\nwant: %q", i, line, tt.expected[i])
				}
			}
		})
	}
}

func TestReadLines_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "nonexistent.txt")

	noop := func(format string, args ...any) {}

	// Should return error for non-existent file
	lines, err := readFileLines(testFile, noop)
	if err == nil {
		t.Errorf("expected error for non-existent file, got nil")
	}
	if lines != nil {
		t.Errorf("expected nil slice for non-existent file, got: %v", lines)
	}
}

// stringListDict tests

func TestStringListDict(t *testing.T) {
	tests := []struct {
		name  string
		input map[string][]string
	}{
		{
			name:  "empty map",
			input: map[string][]string{},
		},
		{
			name: "single key with single value",
			input: map[string][]string{
				"foo.bzl": {"@repo//pkg:file.bzl"},
			},
		},
		{
			name: "single key with multiple values",
			input: map[string][]string{
				"foo.bzl": {
					"@repo1//pkg:file1.bzl",
					"@repo2//pkg:file2.bzl",
				},
			},
		},
		{
			name: "multiple keys sorted alphabetically",
			input: map[string][]string{
				"zebra.bzl": {"@repo//z:z.bzl"},
				"alpha.bzl": {"@repo//a:a.bzl"},
				"beta.bzl":  {"@repo//b:b.bzl"},
			},
		},
		{
			name: "multiple keys with multiple values",
			input: map[string][]string{
				"rules.bzl": {
					"@rules_proto//proto:defs.bzl",
					"@bazel_skylib//lib:paths.bzl",
				},
				"providers.bzl": {
					"@rules_proto//proto:providers.bzl",
				},
			},
		},
		{
			name: "empty values list",
			input: map[string][]string{
				"empty.bzl": {},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringListDict(tt.input)

			// Verify the structure
			if result == nil {
				t.Fatal("result should not be nil")
			}

			if len(tt.input) == 0 && len(result.List) != 0 {
				t.Errorf("empty input should produce empty dict, got %d entries", len(result.List))
			}

			if len(tt.input) != 0 && len(result.List) != len(tt.input) {
				t.Errorf("expected %d dict entries, got %d", len(tt.input), len(result.List))
			}

			// Verify each key-value pair
			for _, kv := range result.List {
				keyExpr, ok := kv.Key.(*build.StringExpr)
				if !ok {
					t.Errorf("key should be StringExpr, got %T", kv.Key)
					continue
				}

				listExpr, ok := kv.Value.(*build.ListExpr)
				if !ok {
					t.Errorf("value should be ListExpr, got %T", kv.Value)
					continue
				}

				// Verify all list elements are StringExpr
				for i, elem := range listExpr.List {
					if _, ok := elem.(*build.StringExpr); !ok {
						t.Errorf("list element %d should be StringExpr, got %T", i, elem)
					}
				}

				// Verify the values match
				expectedValues, exists := tt.input[keyExpr.Value]
				if !exists {
					t.Errorf("unexpected key %q in result", keyExpr.Value)
					continue
				}

				if len(listExpr.List) != len(expectedValues) {
					t.Errorf("key %q: expected %d values, got %d", keyExpr.Value, len(expectedValues), len(listExpr.List))
				}
			}
		})
	}
}

// Exercise both library resolvers with the empty RepoName used when Gazelle
// runs over an extracted archive, as well as named standalone repositories.
func TestLibrarySpecialPackages(t *testing.T) {
	for _, repository := range []struct{ name, indexed, canonical, prefix string }{
		{"archive", "", "extension+fixture", "@@extension+fixture"},
		{"named", "fixture", "", "@fixture"},
		{"main", "", "", ""},
	} {
		for _, library := range []struct{ kind, member, attr string }{
			{starlarkPackageLibraryKind, starlarkPackageKind, "packages"},
			{starlarkModuleLibraryKind, starlarkModuleKind, "modules"},
		} {
			t.Run(repository.name+"/"+library.kind, func(t *testing.T) {
				c := config.New()
				c.RepoName = repository.indexed
				ext := &starlarkRepositoryLang{roots: []string{""}}
				fs := flag.NewFlagSet("test", flag.ContinueOnError)
				ext.RegisterFlags(fs, "fix", c)
				if err := fs.Parse([]string{"-starlarkrepository_canonical_repo_name", repository.canonical}); err != nil {
					t.Fatal(err)
				}
				ix := resolve.NewRuleIndex(func(*rule.Rule, string) resolve.Resolver { return ext })
				for _, pkg := range []string{"", "conditions", "visibility", "ordinary", "conditions/subdir"} {
					member := rule.NewRule(library.member, "member")
					member.SetAttr("src", "defs.bzl")
					ix.AddRule(c, member, &rule.File{Pkg: pkg})
				}
				ix.Finish()
				r := rule.NewRule(library.kind, "library")
				ext.Resolve(c, ix, nil, r, nil, label.New(c.RepoName, "", r.Name()))
				want := []string{":member", "//conditions/subdir:member", "//ordinary:member", repository.prefix + "//conditions:member", repository.prefix + "//visibility:member"}
				sort.Strings(want)
				if got := r.AttrStrings(library.attr); !reflect.DeepEqual(got, want) {
					t.Fatalf("%s = %v, want %v", library.attr, got, want)
				}
			})
		}
	}
}

func TestLibraryDependencyLabel(t *testing.T) {
	ext := &starlarkRepositoryLang{canonicalRepoName: "extension+fixture"}
	for _, tc := range []struct{ dep, from, want string }{
		{"@fixture//conditions:member", "@fixture//conditions:library", ":member"},
		{"@fixture//visibility:member", "@fixture//visibility:library", ":member"},
		{"@other//conditions:member", "@fixture//:library", "@other//conditions:member"},
		{"@@other+repo//visibility:member", "@fixture//:library", "@@other+repo//visibility:member"},
		{"@fixture//ordinary:member", "@fixture//:library", "//ordinary:member"},
	} {
		t.Run(tc.dep+" from "+tc.from, func(t *testing.T) {
			dep, err := label.Parse(tc.dep)
			if err != nil {
				t.Fatal(err)
			}
			from, err := label.Parse(tc.from)
			if err != nil {
				t.Fatal(err)
			}
			if got := ext.libraryDependencyLabel(dep, from); got != tc.want {
				t.Errorf("got %q, want %q", got, tc.want)
			}
		})
	}
}
