package prost

import "testing"

func TestNeedsCompileWellKnownTypes(t *testing.T) {
	cases := []struct {
		name string
		opts []string
		want bool
	}{
		{"empty", nil, false},
		{"unrelated extern", []string{"extern_path=.foo.bar=::foo_bar"}, false},
		{
			"per-package google.protobuf",
			[]string{"extern_path=.google.protobuf=::google_protobuf__descriptor"},
			true,
		},
		{
			"per-type google.protobuf.Empty",
			[]string{"extern_path=.google.protobuf.Empty=::google_protobuf__empty::Empty"},
			true,
		},
		{
			// `google.protobufabc` is a different package, must not match.
			"prefix-but-not-google.protobuf",
			[]string{"extern_path=.google.protobufabc=::google_protobufabc"},
			false,
		},
	}
	for _, c := range cases {
		if got := needsCompileWellKnownTypes(c.opts, nil); got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, got, c.want)
		}
	}
}
