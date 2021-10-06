package protoc

import (
	"testing"

	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/emicklei/proto"
)

func TestParseRewrite(t *testing.T) {
	for name, tc := range map[string]struct {
		spec string
		in   string
		want string
	}{
		"basic regexp": {
			spec: "a b",
			in:   "a",
			want: "b",
		},
		"canonical example": {
			spec: "google/protobuf/([a-z]+).proto @org_golang_google_protobuf//types/known/${1}pb",
			in:   "google/protobuf/empty.proto",
			want: "@org_golang_google_protobuf//types/known/emptypb",
		},
	} {
		t.Run(name, func(t *testing.T) {
			rw, err := ParseRewrite(tc.spec)
			if err != nil {
				t.Fatal(err)
			}
			got := rw.ReplaceAll(tc.in)
			if got != tc.want {
				t.Errorf("replaceall: want %s, got %s", tc.want, got)
			}
		})
	}
}

func TestResolveLibraryRewrites(t *testing.T) {
	for name, tc := range map[string]struct {
		imports  []string
		rewrites []string
		want     []string
	}{
		"single match": {
			imports:  []string{"google/protobuf/any.proto"},
			rewrites: []string{"google/protobuf/([a-z]+).proto @org_golang_google_protobuf//types/known/${1}pb"},
			want:     []string{"@org_golang_google_protobuf//types/known/anypb"},
		},
		"multiple match": {
			imports: []string{
				"google/protobuf/any.proto",
				"google/protobuf/empty.proto",
				"not/matched.proto",
			},
			rewrites: []string{"google/protobuf/([a-z]+).proto @org_golang_google_protobuf//types/known/${1}pb"},
			want: []string{
				"@org_golang_google_protobuf//types/known/anypb",
				"@org_golang_google_protobuf//types/known/emptypb",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			file := NewFile("fake/rel/dir", "foo")
			for _, v := range tc.imports {
				file.imports = append(file.imports, proto.Import{Filename: v})
			}
			lib := NewOtherProtoLibrary(
				rule.EmptyFile("", "fake/rel/dir"),
				rule.NewRule("proto_library", "foo_proto_library"),
				file)
			rewrites := make([]Rewrite, len(tc.rewrites))
			for i, d := range tc.rewrites {
				rw, err := ParseRewrite(d)
				if err != nil {
					t.Fatal("parse", d, err)
				}
				rewrites[i] = *rw
			}
			got := ResolveLibraryRewrites(rewrites, lib)
			if len(tc.want) != len(got) {
				t.Fatalf("want %d, got %d (%v -vs- %v) ", len(tc.want), len(got), tc.want, got)
			}
			for i, actual := range got {
				expected := tc.want[i]
				if actual != expected {
					t.Errorf("item %d: want %s, got %s", i, expected, actual)
				}
			}
		})
	}
}
