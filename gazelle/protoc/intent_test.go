package protoc

import "testing"

func TestParseIntent(t *testing.T) {
	for name, tc := range map[string]struct {
		in        string
		wantValue string
		wantNeg   bool
	}{
		"": {
			in:        "",
			wantValue: "",
			wantNeg:   false,
		},
		"bare": {
			in:        "foo",
			wantValue: "foo",
			wantNeg:   false,
		},
		"+": {
			in:        "+foo",
			wantValue: "foo",
			wantNeg:   false,
		},
		"-": {
			in:        "-foo",
			wantValue: "foo",
			wantNeg:   true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			intent := parseIntent(tc.in)
			if tc.wantValue != intent.Value {
				t.Errorf("value: want %s, got %s", tc.wantValue, intent.Value)
			}
			if tc.wantNeg != intent.Negative {
				t.Errorf("value: want %t, got %t", tc.wantNeg, intent.Negative)
			}
		})
	}
}
