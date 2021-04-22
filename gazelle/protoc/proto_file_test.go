package protoc

import (
	"strings"
	"testing"
)

func mustParseTestProtoFile(t *testing.T, in string) *ProtoFile {
	f := &ProtoFile{}
	if err := f.parseReader(strings.NewReader(in)); err != nil {
		t.Fatalf("mustTestProtoFile: %v", err)
	}
	return f
}

func TestHas(t *testing.T) {
	type hasTestCase struct {
		in            string
		hasMessages   bool
		hasServices   bool
		hasEnumOption string
	}
	tests := map[string]*hasTestCase{
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			f := mustParseTestProtoFile(t, tc.in)
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
		})
	}
}
