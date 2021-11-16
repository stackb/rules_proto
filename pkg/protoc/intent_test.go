package protoc

import "testing"

func TestParseIntent(t *testing.T) {
	for name, tc := range map[string]struct {
		in        string
		wantValue string
		want      bool
	}{
		"": {
			in:        "",
			wantValue: "",
			want:      true,
		},
		"bare": {
			in:        "foo",
			wantValue: "foo",
			want:      true,
		},
		"+": {
			in:        "+foo",
			wantValue: "foo",
			want:      true,
		},
		"-": {
			in:        "-foo",
			wantValue: "foo",
			want:      false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			intent := parseIntent(tc.in)
			if tc.wantValue != intent.Value {
				t.Errorf("value: want %s, got %s", tc.wantValue, intent.Value)
			}
			if tc.want != intent.Want {
				t.Errorf("value: want %t, got %t", tc.want, intent.Want)
			}
		})
	}
}
