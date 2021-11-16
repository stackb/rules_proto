// gencopy is a utility program that copies bazel outputs back into the
// workspace source tree.  Ideally, you don't have any generated files committed
// to VCS, but sometimes you do.
//
package main

import (
	"os"
	"testing"
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
			want: os.FileMode(0644),
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
