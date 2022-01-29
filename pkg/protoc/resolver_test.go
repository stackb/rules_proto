package protoc

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/google/go-cmp/cmp"
)

func TestLoadResolver(t *testing.T) {
	for name, tc := range map[string]struct {
		in    string
		known map[string]importLabels
	}{
		"empty string": {
			in:    "",
			known: map[string]importLabels{},
		},
		"comment": {
			in:    "# ignored",
			known: map[string]importLabels{},
		},
		"proto resolve": {
			in: "proto,proto,google/protobuf/any.proto,@com_google_protobuf//:any_proto",
			known: map[string]importLabels{
				"proto proto": map[string][]label.Label{
					"google/protobuf/any.proto": {label.New("com_google_protobuf", "", "any_proto")},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			resolver := &resolver{
				options: &ImportResolverOptions{
					Debug:  false,
					Printf: t.Logf,
				},
				known: make(map[string]importLabels),
			}
			if err := resolver.Load(strings.NewReader(tc.in)); err != nil {
				t.Fatal(err)
			}
			if diff := cmp.Diff(tc.known, resolver.known); diff != "" {
				t.Error("unexpected diff:", diff)
			}
		})
	}
}

func TestSaveResolver(t *testing.T) {
	for name, tc := range map[string]struct {
		known    map[string]importLabels
		repoName string
		out      string
	}{
		"empty string": {
			known: map[string]importLabels{},
			out:   "",
		},
		"does not write external labels": {
			// the resolver generally loads resolves from other csv files, so we
			// don't want to "transitively" emit them.  Saved output should onlt
			// reflect the current workspace.
			known: map[string]importLabels{
				"proto proto": map[string][]label.Label{
					"google/protobuf/any.proto": {label.New("com_google_protobuf", "", "any_proto")},
				},
			},
			out: "",
		},
		"rewrites labels with repoName": {
			repoName: "com_google_protobuf",
			known: map[string]importLabels{
				"proto proto": map[string][]label.Label{
					"google/protobuf/any.proto": {label.New("", "", "any_proto")},
				},
			},
			out: "proto,proto,google/protobuf/any.proto,@com_google_protobuf//:any_proto\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			resolver := &resolver{
				options: &ImportResolverOptions{
					Debug:  false,
					Printf: t.Logf,
				},
				known: tc.known,
			}
			var out bytes.Buffer
			resolver.Save(&out, tc.repoName)
			if diff := cmp.Diff(tc.out, out.String()); diff != "" {
				t.Error("unexpected diff:", diff)
			}
		})
	}
}

func TestProvide(t *testing.T) {
	for name, tc := range map[string]struct {
		lang, impLang, imp string
		from               label.Label
		known              map[string]importLabels
	}{
		"empty case": {
			known: map[string]importLabels{
				" ": map[string][]label.Label{
					"": {label.NoLabel},
				},
			},
		},
		"typical usage": {
			lang:    "proto",
			impLang: "proto",
			imp:     "google/protobuf/any.proto",
			from:    label.New("com_google_protobuf", "", "any_proto"),
			known: map[string]importLabels{
				"proto proto": map[string][]label.Label{
					"google/protobuf/any.proto": {label.New("com_google_protobuf", "", "any_proto")},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			resolver := &resolver{
				options: &ImportResolverOptions{
					Debug:  false,
					Printf: t.Logf,
				},
				known: make(map[string]importLabels),
			}
			resolver.Provide(tc.lang, tc.impLang, tc.imp, tc.from)
			if diff := cmp.Diff(tc.known, resolver.known); diff != "" {
				t.Error("unexpected diff:", diff)
			}
		})
	}
}

func TestResolve(t *testing.T) {
	for name, tc := range map[string]struct {
		lang, impLang, imp string
		want               []resolve.FindResult
		known              map[string]importLabels
	}{
		"empty case - matches a single empty result": {
			known: map[string]importLabels{
				" ": map[string][]label.Label{
					"": {label.NoLabel},
				},
			},
			want: []resolve.FindResult{{}},
		},
		"typical usage": {
			lang:    "proto",
			impLang: "proto",
			imp:     "google/protobuf/any.proto",
			want: []resolve.FindResult{
				{
					Label: label.New("com_google_protobuf", "", "any_proto"),
				},
			},
			known: map[string]importLabels{
				"proto proto": map[string][]label.Label{
					"google/protobuf/any.proto": {label.New("com_google_protobuf", "", "any_proto")},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			resolver := &resolver{
				options: &ImportResolverOptions{
					Debug:  false,
					Printf: t.Logf,
				},
				known: tc.known,
			}
			got := resolver.Resolve(tc.lang, tc.impLang, tc.imp)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("Resolve() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
