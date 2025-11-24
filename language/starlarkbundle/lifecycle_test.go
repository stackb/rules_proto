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
	"context"
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/config"
)

func TestLoggerLifecycle(t *testing.T) {
	// Create a temporary directory for the log file
	tmpDir := t.TempDir()
	logFile := filepath.Join(tmpDir, "test.log")

	// Create extension and register flags
	ext := NewLanguage().(*starlarkBundleLang)
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	c := &config.Config{}
	ext.RegisterFlags(fs, "update", c)

	// Set the log file flag
	if err := fs.Set("starlark_bundle_log", logFile); err != nil {
		t.Fatalf("failed to set flag: %v", err)
	}

	// Initialize logger via CheckFlags
	if err := ext.CheckFlags(fs, c); err != nil {
		t.Fatalf("CheckFlags failed: %v", err)
	}

	// Verify logger is initialized
	if ext.logger == nil {
		t.Fatal("logger should be initialized")
	}
	if ext.logWriter == nil {
		t.Fatal("logWriter should be initialized")
	}

	// Write some log messages
	ext.logf("test message 1")
	ext.logf("test message 2: %s", "with formatting")

	// Call lifecycle methods
	ctx := context.Background()
	ext.Before(ctx)
	ext.DoneGeneratingRules()
	ext.AfterResolvingDeps(ctx)

	// Verify log writer is closed
	if ext.logWriter != nil {
		t.Error("logWriter should be nil after AfterResolvingDeps")
	}

	// Read the log file and verify content
	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("failed to read log file: %v", err)
	}

	logContent := string(content)
	if !strings.Contains(logContent, "test message 1") {
		t.Error("log file should contain 'test message 1'")
	}
	if !strings.Contains(logContent, "test message 2: with formatting") {
		t.Error("log file should contain 'test message 2: with formatting'")
	}
}

func TestLoggerWithoutFile(t *testing.T) {
	// Create extension without log file
	ext := NewLanguage().(*starlarkBundleLang)
	fs := flag.NewFlagSet("test", flag.ContinueOnError)
	c := &config.Config{}
	ext.RegisterFlags(fs, "update", c)

	// Initialize logger via CheckFlags (no log file set)
	if err := ext.CheckFlags(fs, c); err != nil {
		t.Fatalf("CheckFlags failed: %v", err)
	}

	// Verify logger is initialized but logWriter is nil
	if ext.logger == nil {
		t.Fatal("logger should be initialized")
	}
	if ext.logWriter != nil {
		t.Error("logWriter should be nil when no log file is specified")
	}

	// Writing logs should not panic
	ext.logf("this should not panic")

	// Lifecycle methods should not panic
	ctx := context.Background()
	ext.AfterResolvingDeps(ctx)
}
