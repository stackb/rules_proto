package protoc

import (
	"testing"
)

func TestStripRel(t *testing.T) {
	for name, tc := range map[string]struct {
		rel  string
		in   string
		want string
	}{
		"empty": {
			rel:  "",
			in:   "",
			want: "",
		},
		"match": {
			rel:  "proto",
			in:   "proto/foo.txt",
			want: "foo.txt",
		},
		"no-match": {
			rel:  "proto",
			in:   "foo/proto.txt",
			want: "foo/proto.txt",
		},
	} {
		t.Run(name, func(t *testing.T) {
			got := StripRel(tc.rel, tc.in)
			if got != tc.want {
				t.Errorf("striprel: want %s, got %s", tc.want, got)
			}
		})
	}
}
