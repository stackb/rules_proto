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
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveAbsolutePathsFromBazelignore(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		removed  bool
	}{
		{
			name: "removes absolute paths",
			input: `.build/
/.github/
.swiftpm/
/.vscode/
`,
			expected: `.build/
.swiftpm/
`,
			removed: true,
		},
		{
			name: "keeps relative paths",
			input: `.build/
.github/
.swiftpm/
.vscode/
`,
			expected: `.build/
.github/
.swiftpm/
.vscode/
`,
			removed: false,
		},
		{
			name: "handles empty lines",
			input: `.build/

/.github/

.swiftpm/
`,
			expected: `.build/


.swiftpm/
`,
			removed: true,
		},
		{
			name: "handles comments and absolute paths",
			input: `# Comment
.build/
/absolute/path
relative/path
/another/absolute
`,
			expected: `# Comment
.build/
relative/path
`,
			removed: true,
		},
		{
			name:     "empty file",
			input:    "",
			expected: "",
			removed:  false,
		},
		{
			name: "only absolute paths",
			input: `/path1
/path2
/path3
`,
			expected: ``,
			removed: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temp directory
			tmpDir := t.TempDir()
			bazelignorePath := filepath.Join(tmpDir, ".bazelignore")

			// Write input content
			if err := os.WriteFile(bazelignorePath, []byte(tt.input), 0644); err != nil {
				t.Fatalf("failed to write test file: %v", err)
			}

			// Create extension
			ext := NewLanguage().(*starlarkBundleLang)
			ext.logger = nil // Disable logging for tests

			// Run the function
			err := ext.removeAbsolutePathsFromBazelignore(bazelignorePath)
			if err != nil {
				t.Fatalf("removeAbsolutePathsFromBazelignore failed: %v", err)
			}

			// Read result
			result, err := os.ReadFile(bazelignorePath)
			if err != nil {
				t.Fatalf("failed to read result file: %v", err)
			}

			if string(result) != tt.expected {
				t.Errorf("Content mismatch:\ngot:\n%q\n\nwant:\n%q", string(result), tt.expected)
			}
		})
	}
}

func TestRemoveAbsolutePathsFromBazelignore_NonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	bazelignorePath := filepath.Join(tmpDir, ".bazelignore")

	ext := NewLanguage().(*starlarkBundleLang)
	ext.logger = nil

	// Should not error on non-existent file
	err := ext.removeAbsolutePathsFromBazelignore(bazelignorePath)
	if err != nil {
		t.Errorf("expected no error for non-existent file, got: %v", err)
	}
}
