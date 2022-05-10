package rules_scala

import (
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
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
