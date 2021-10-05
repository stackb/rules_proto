package protoc

import "testing"

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
