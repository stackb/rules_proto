package protoc

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func mustParseTestFile(t *testing.T, in string) *File {
	f := &File{}
	if err := f.ParseReader(strings.NewReader(in)); err != nil {
		t.Fatalf("mustTestFile: %v", err)
	}
	return f
}

func TestHas(t *testing.T) {
	tests := map[string]struct {
		in            string
		hasMessages   bool
		hasServices   bool
		hasEnumOption string
		hasRPCOption  string
	}{
		"empty file": {},
		"has services": {
			in: `
syntax = "proto3";

service Greeter {
	rpc Greet(GreetRequest) returns (GreetResponse);
}
`,
			hasServices: true,
		},
		"has enum option": {
			in: `
syntax = "proto3";
import "google/api/visibility.proto";
enum MyEnum {
    UNKNOWN = 0;
	PRIVATE = 1 [(google.api.value_visibility).restriction = "HIDDEN"];
}
`,
			hasEnumOption: "(google.api.value_visibility).restriction",
		},
		"has rpc option": {
			in: `
syntax = "proto3";
import "google/api/annotations.proto";

service Greeter {
	rpc Greet(GreetRequest) returns (GreetResponse) {
		option (google.api.http) = {
			get: "/greet"
		};
	}
}
`,
			hasServices:  true,
			hasRPCOption: "(google.api.http)",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := mustParseTestFile(t, tc.in)
			if tc.hasMessages != f.HasMessages() {
				t.Errorf("hasMessages: want %t, got %t", tc.hasMessages, f.HasMessages())
			}
			if tc.hasServices != f.HasServices() {
				t.Errorf("hasServices: want %t, got %t", tc.hasServices, f.HasServices())
			}
			if tc.hasEnumOption != "" && !f.HasEnumOption(tc.hasEnumOption) {
				t.Errorf("hasEnumOption: expected %s",
					tc.hasEnumOption)
			}
			if tc.hasRPCOption != "" && !f.HasRPCOption(tc.hasRPCOption) {
				t.Errorf("hasRPCOption: expected %s",
					tc.hasRPCOption)
			}
		})
	}
}

func TestRelativeFileNameWithExtensions(t *testing.T) {
	tests := map[string]struct {
		dir  string
		name string
		rel  string
		exts []string
		want []string
	}{
		"empty": {
			want: []string{},
		},
		"single": {
			name: "a",
			rel:  "proto",
			exts: []string{".cc"},
			want: []string{
				"proto/a.cc",
			},
		},
		"multiple": {
			name: "a",
			rel:  "proto",
			exts: []string{".cc", ".h"},
			want: []string{
				"proto/a.cc",
				"proto/a.h",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fn := RelativeFileNameWithExtensions(tc.rel, tc.exts...)
			got := fn(&File{
				Dir:  tc.dir,
				Name: tc.name,
			})
			assert.Equal(t, got, tc.want, "generated filenames")
		})
	}
}

func TestImportPrefixRelativeFileNameWithExtensions(t *testing.T) {
	tests := map[string]struct {
		stripImportPrefix string
		dir               string
		name              string
		rel               string
		exts              []string
		want              []string
	}{
		"empty": {
			want: []string{},
		},
		"single": {
			name: "a",
			rel:  "proto",
			exts: []string{".cc"},
			want: []string{
				"proto/a.cc",
			},
		},
		"multiple": {
			name: "a",
			rel:  "proto",
			exts: []string{".cc", ".h"},
			want: []string{
				"proto/a.cc",
				"proto/a.h",
			},
		},
		"strip": {
			stripImportPrefix: "foo/bar",
			name:              "a",
			rel:               "foo/bar/baz",
			exts:              []string{".cc", ".h"},
			want: []string{
				"baz/a.cc",
				"baz/a.h",
			},
		},
		"strip-abs": {
			stripImportPrefix: "/foo/bar",
			name:              "a",
			rel:               "foo/bar/baz",
			exts:              []string{".cc", ".h"},
			want: []string{
				"baz/a.cc",
				"baz/a.h",
			},
		},
		"strip-abs-with-trailing-slash": {
			stripImportPrefix: "/foo/bar/",
			name:              "a",
			rel:               "foo/bar/baz",
			exts:              []string{".cc", ".h"},
			want: []string{
				"baz/a.cc",
				"baz/a.h",
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			fn := ImportPrefixRelativeFileNameWithExtensions(tc.stripImportPrefix, tc.rel, tc.exts...)
			got := fn(&File{
				Dir:  tc.dir,
				Name: tc.name,
			})
			assert.Equal(t, got, tc.want, "generated filenames")
		})
	}
}
