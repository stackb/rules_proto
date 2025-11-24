/* Copyright 2020 The Bazel Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package starlarkbundle

import (
	"strings"
	"testing"

	"github.com/bazelbuild/buildtools/build"
)

func TestRenameFailToPrint(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		want     string
		modified bool
	}{
		{
			name: "single fail call",
			input: `fail("error message")
`,
			want: `print("error message")
`,
			modified: true,
		},
		{
			name: "multiple fail calls",
			input: `if condition:
    fail("error 1")
else:
    fail("error 2")
`,
			want: `if condition:
    print("error 1")
else:
    print("error 2")
`,
			modified: true,
		},
		{
			name: "fail with multiple arguments",
			input: `fail("error:", variable, "more text")
`,
			want: `print("error:", variable, "more text")
`,
			modified: true,
		},
		{
			name: "no fail calls",
			input: `print("hello")
x = 1
`,
			want: `print("hello")
x = 1
`,
			modified: false,
		},
		{
			name: "mixed fail and print calls",
			input: `print("info")
fail("error")
print("more info")
`,
			want: `print("info")
print("error")
print("more info")
`,
			modified: true,
		},
		{
			name: "fail in function definition",
			input: `def my_function(param):
    if not param:
        fail("param is required")
    return param
`,
			want: `def my_function(param):
    if not param:
        print("param is required")
    return param
`,
			modified: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ast, err := build.ParseBzl("test.bzl", []byte(tt.input))
			if err != nil {
				t.Fatalf("failed to parse input: %v", err)
			}

			modified := renameFailToPrint(ast)

			if modified != tt.modified {
				t.Errorf("renameFailToPrint() modified = %v, want %v", modified, tt.modified)
			}

			got := string(build.Format(ast))
			if strings.TrimSpace(got) != strings.TrimSpace(tt.want) {
				t.Errorf("renameFailToPrint() output mismatch\ngot:\n%s\n\nwant:\n%s", got, tt.want)
			}
		})
	}
}
