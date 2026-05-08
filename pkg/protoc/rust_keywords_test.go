package protoc

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRustKeywordEscapeMappings(t *testing.T) {
	for name, tc := range map[string]struct {
		pkg     string
		outputs []string
		want    map[string]string
	}{
		"empty package": {
			pkg:     "",
			outputs: []string{"foo.rs"},
			want:    nil,
		},
		"empty outputs": {
			pkg:     "google.type",
			outputs: nil,
			want:    nil,
		},
		"no keywords": {
			pkg:     "google.api",
			outputs: []string{"google.api.rs", "google.api.serde.rs"},
			want:    nil,
		},
		"type keyword": {
			pkg:     "google.type",
			outputs: []string{"google.type.rs", "google.type.serde.rs"},
			want: map[string]string{
				"google.type.rs":       "google/r#type/google.type.rs",
				"google.type.serde.rs": "google/r#type/google.type.serde.rs",
			},
		},
		"keyword at start": {
			pkg:     "type.example",
			outputs: []string{"type.example.rs"},
			want: map[string]string{
				"type.example.rs": "r#type/example/type.example.rs",
			},
		},
		"multiple keywords": {
			pkg:     "self.type",
			outputs: []string{"self.type.rs"},
			want: map[string]string{
				"self.type.rs": "r#self/r#type/self.type.rs",
			},
		},
		"single segment keyword": {
			pkg:     "type",
			outputs: []string{"type.rs"},
			want: map[string]string{
				"type.rs": "r#type/type.rs",
			},
		},
		"single segment no keyword": {
			pkg:     "example",
			outputs: []string{"example.rs"},
			want:    nil,
		},
		"full path outputs": {
			pkg:     "google.type",
			outputs: []string{"google/type/google.type.rs", "google/type/google.type.serde.rs"},
			want: map[string]string{
				"google.type.rs":       "google/r#type/google.type.rs",
				"google.type.serde.rs": "google/r#type/google.type.serde.rs",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := RustKeywordEscapeMappings(tc.pkg, tc.outputs)
			if diff := cmp.Diff(tc.want, got); diff != "" {
				t.Errorf("(-want +got):\n%s", diff)
			}
		})
	}
}
