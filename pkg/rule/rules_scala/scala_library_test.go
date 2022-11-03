package rules_scala

import (
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/google/go-cmp/cmp"
	"github.com/stackb/rules_proto/pkg/protoc"
)

// TestGetJavaPackageOption shows that an import named in (scalapb.options) works as expected.
func TestGetJavaPackageOption(t *testing.T) {
	for name, tc := range map[string]struct {
		in   string
		want string
	}{
		"degenerate case": {},
		"with go_package": {
			in: `syntax="proto3"; option go_package="com.foo";`,
		},
		"with java_package": {
			in:   `syntax="proto3"; option java_package="com.foo";`,
			want: "com.foo",
		},
	} {
		t.Run(name, func(t *testing.T) {
			file := protoc.NewFile("", "test.proto")
			if err := file.ParseReader(strings.NewReader(tc.in)); err != nil {
				t.Fatal("parse file:", err)
			}
			got, ok := javaPackageOption(file.Options())
			if ok {
				if diff := cmp.Diff(tc.want, got); diff != "" {
					t.Errorf("TestGetScalaImports() mismatch (-want +got):\n%s", diff)
				}
			} else {
				if tc.want != "" {
					t.Errorf("TestGetScalaImports() unexpected miss: %v", tc)
				}
			}
		})
	}
}

// TestGetScalapbImports shows that an import named in (scalapb.options) works as expected.
func TestGetScalapbImports(t *testing.T) {
	for name, tc := range map[string]struct {
		// in is a mapping of source filename to content
		in   map[string]string
		want []string
	}{
		"degenerate case": {
			want: []string{},
		},
		"without imports": {
			in: map[string]string{
				"foo.proto": `syntax = "proto3";
message Thing {}`,
			},
			want: []string{},
		},
		"with scalapb import": {
			in: map[string]string{
				"foo.proto": `syntax = "proto3";
import "scalapb/scalapb.proto";

option (scalapb.options) = {
	import: "corp.common.utils.WithORM"
};`,
			},
			want: []string{"corp.common.utils.WithORM"},
		},
		"with field type": {
			in: map[string]string{
				"foo.proto": `
syntax = "proto2";

import "thirdparty/protobuf/scalapb/scalapb.proto";

message TraderId {
	required int32 trader_id = 1 [(scalapb.field).type = "corp.common.utils.TraderId"];
}

message TeamId {
	required int32 team_id = 1 [(scalapb.field).type = "corp.common.utils.TeamId"];
}				
`,
			},
			want: []string{"corp.common.utils.TeamId", "corp.common.utils.TraderId"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			files := make([]*protoc.File, len(tc.in))
			i := 0
			for name, content := range tc.in {
				file := protoc.NewFile("", name)
				if err := file.ParseReader(strings.NewReader(content)); err != nil {
					t.Fatal("parse file:", name, err)
				}
				files[i] = file
				i++
			}
			got := getScalapbImports(files)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("TestGetScalaImports() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

// TestProvideScalaImports shows the imports provided.
func TestProvideScalaImports(t *testing.T) {
	for name, tc := range map[string]struct {
		// in is a mapping of source filename to content
		in map[string]string
		// options is a mapping of protoc options
		options map[string]bool
		want    []resolve.ImportSpec
	}{
		"degenerate case": {},
		"message": {
			in: map[string]string{
				"foo.proto": `syntax = "proto3";
message Thing {}`,
			},
			want: []resolve.ImportSpec{
				{Lang: "scala", Imp: "Thing"},
				{Lang: "scala", Imp: "ThingProto"},
			},
		},
		"service": {
			in: map[string]string{
				"foo.proto": `syntax = "proto3";
service Thinger {}`,
			},
			want: []resolve.ImportSpec{
				{Lang: "scala", Imp: "Thinger"},
				{Lang: "scala", Imp: "ThingerGrpc"},
				{Lang: "scala", Imp: "ThingerProto"},
				{Lang: "scala", Imp: "ThingerClient"},
				{Lang: "scala", Imp: "ThingerHandler"},
				{Lang: "scala", Imp: "ThingerServer"},
				{Lang: "scala", Imp: "ThingerPowerApi"},
				{Lang: "scala", Imp: "ThingerPowerApiHandler"},
				{Lang: "scala", Imp: "ThingerClientPowerApi"},
			},
		},
		"enum": {
			in: map[string]string{
				"foo.proto": `syntax = "proto3";
enum Things {}`,
			},
			want: []resolve.ImportSpec{
				{Lang: "scala", Imp: "Things"},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			files := make([]*protoc.File, len(tc.in))
			i := 0
			for name, content := range tc.in {
				file := protoc.NewFile("", name)
				if err := file.ParseReader(strings.NewReader(content)); err != nil {
					t.Fatal("parse file:", name, err)
				}
				files[i] = file
				i++
			}
			resolver := &fakeImportResolver{}
			from := label.New("repo", "dir", "name")

			provideScalaImports(files, resolver, from, tc.options)
			if diff := cmp.Diff(tc.want, resolver.got); diff != "" {
				t.Errorf("TestGetScalaImports() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

type fakeImportResolver struct {
	got []resolve.ImportSpec
}

func (r *fakeImportResolver) Imports(lang, impLang string, visitor func(imp string, location []label.Label) bool) {
	panic("not implemented")
}

func (r *fakeImportResolver) Resolve(lang, impLang, imp string) []resolve.FindResult {
	panic("not implemented")
}

func (r *fakeImportResolver) Provide(lang, impLang, imp string, from label.Label) {
	r.got = append(r.got, resolve.ImportSpec{Imp: imp, Lang: impLang})
}

func TestScalaLibraryOptionsNoResolve(t *testing.T) {
	for name, tc := range map[string]struct {
		args    []string
		imports []string
		want    []string
	}{
		"degenerate case": {},
		"prototypical": {
			args:    []string{"--noresolve=scalapb/scalapb.proto"},
			imports: []string{"scalapb/scalapb.proto", "google/protobuf/any.proto"},
			want:    []string{"google/protobuf/any.proto"},
		},
		"csv": {
			args:    []string{"--noresolve=a.proto,b.proto"},
			imports: []string{"a.proto", "b.proto"},
			want:    nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			options := parseScalaLibraryOptions("proto_scala_library", tc.args)
			got := options.filterImports(tc.imports)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}

func TestScalaLibraryOptionsNoOutput(t *testing.T) {
	for name, tc := range map[string]struct {
		args    []string
		outputs []string
		want    []string
	}{
		"degenerate case": {},
		"prototypical": {
			args:    []string{"--exclude=package_scala.srcjar"},
			outputs: []string{"package_scala.srcjar"},
			want:    nil,
		},
		"csv": {
			args:    []string{"--exclude=a.srcjar,b.srcjar"},
			outputs: []string{"a.srcjar", "b.srcjar"},
			want:    nil,
		},
		"pattern": {
			args:    []string{"--exclude=**/*.srcjar"},
			outputs: []string{"a.srcjar", "lib/b.srcjar", "lib/c.jar"},
			want:    []string{"lib/c.jar"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			options := parseScalaLibraryOptions("proto_scala_library", tc.args)
			got := options.filterOutputs(tc.outputs)

			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}

func TestResolveScalaDeps(t *testing.T) {
	for name, tc := range map[string]struct {
		overrideFn     findRuleWithOverride
		byImportFn     findRulesByImportWithConfig
		r              *rule.Rule
		from           label.Label
		unresolvedDeps map[string]error
		wantUnresolved map[string]error
		wantDeps       []string
	}{
		"degenerate case": {
			overrideFn: func(c *config.Config, imp resolve.ImportSpec, lang string) (label.Label, bool) {
				return label.NoLabel, false
			},
			byImportFn: func(c *config.Config, imp resolve.ImportSpec, lang string) []resolve.FindResult {
				return nil
			},
			wantUnresolved: map[string]error{},
		},
		"resolve from cross-resolver": {
			from: label.New("", "proto", "foo_proto_scala_library"),
			overrideFn: func(c *config.Config, imp resolve.ImportSpec, lang string) (label.Label, bool) {
				return label.NoLabel, false
			},
			byImportFn: func(c *config.Config, imp resolve.ImportSpec, lang string) []resolve.FindResult {
				if lang == "scala" && imp.Imp == "foo.bar.baz.mapper" {
					return []resolve.FindResult{{Label: label.New("", "mapper", "scala_lib")}}
				}
				return nil
			},
			unresolvedDeps: map[string]error{
				"foo.bar.baz.mapper": protoc.ErrNoLabel,
			},
			wantUnresolved: map[string]error{},
			wantDeps:       []string{"//mapper:scala_lib"},
		},
		"resolve from overrideFn": {
			from: label.New("", "proto", "foo_proto_scala_library"),
			overrideFn: func(c *config.Config, imp resolve.ImportSpec, lang string) (label.Label, bool) {
				if imp.Lang == "scala" && imp.Imp == "foo.bar.baz.mapper" {
					return label.New("", "mapper", "scala_lib"), true
				}
				return label.NoLabel, false
			},
			byImportFn: func(c *config.Config, imp resolve.ImportSpec, lang string) []resolve.FindResult {
				return nil
			},
			unresolvedDeps: map[string]error{
				"foo.bar.baz.mapper": protoc.ErrNoLabel,
			},
			wantUnresolved: map[string]error{},
			wantDeps:       []string{"//mapper:scala_lib"},
		},
		"does not resolve self-label": {
			from: label.New("", "proto", "foo_proto_scala_library"),
			overrideFn: func(c *config.Config, imp resolve.ImportSpec, lang string) (label.Label, bool) {
				if imp.Lang == "scala" && imp.Imp == "foo.bar.baz.mapper" {
					return label.New("", "proto", "foo_proto_scala_library"), true
				}
				return label.NoLabel, false
			},
			byImportFn: func(c *config.Config, imp resolve.ImportSpec, lang string) []resolve.FindResult {
				return nil
			},
			unresolvedDeps: map[string]error{
				"foo.bar.baz.mapper": protoc.ErrNoLabel,
			},
			wantUnresolved: map[string]error{
				"foo.bar.baz.mapper": protoc.ErrNoLabel,
			},
			wantDeps: nil,
		},
	} {
		t.Run(name, func(t *testing.T) {
			c := &config.Config{}
			r := rule.NewRule("proto_scala_library", "bar_proto_scala_library")

			got := make(map[string]error)
			for k, v := range tc.unresolvedDeps {
				got[k] = v
			}
			resolveScalaDeps(tc.overrideFn, tc.byImportFn, c, r, got, tc.from)

			gotDeps := r.AttrStrings("deps")

			if diff := cmp.Diff(tc.wantDeps, gotDeps); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}

type fakeCrossResolver struct {
	result []resolve.FindResult
}

func (cr *fakeCrossResolver) CrossResolve(c *config.Config, ix *resolve.RuleIndex, imp resolve.ImportSpec, lang string) []resolve.FindResult {
	return cr.result
}
