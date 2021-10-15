package protoc

import (
	"bytes"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/label"
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
					"google/protobuf/any.proto": []label.Label{label.New("com_google_protobuf", "", "any_proto")},
				},
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			resolver := &resolver{
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
					"google/protobuf/any.proto": []label.Label{label.New("com_google_protobuf", "", "any_proto")},
				},
			},
			out: "",
		},
		"rewrites labels with repoName": {
			repoName: "com_google_protobuf",
			known: map[string]importLabels{
				"proto proto": map[string][]label.Label{
					"google/protobuf/any.proto": []label.Label{label.New("", "", "any_proto")},
				},
			},
			out: "proto,proto,google/protobuf/any.proto,@com_google_protobuf//:any_proto\n",
		},
	} {
		t.Run(name, func(t *testing.T) {
			resolver := &resolver{tc.known}
			var out bytes.Buffer
			resolver.Save(&out, tc.repoName)
			if diff := cmp.Diff(tc.out, out.String()); diff != "" {
				t.Error("unexpected diff:", diff)
			}
		})
	}
}
