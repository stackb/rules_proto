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
)

// Before implements part of the language.LifecycleManager interface.
func (ext *starlarkBundleLang) Before(context.Context) {
	ext.logf("Lifecycle: Before() called")
}

// DoneGeneratingRules implements part of the language.FinishableLanguage interface.
func (ext *starlarkBundleLang) DoneGeneratingRules() {
	ext.logf("Lifecycle: DoneGeneratingRules() called")
}

// AfterResolvingDeps implements part of the language.LifecycleManager interface.
// This is where we flush and close the log file.
func (ext *starlarkBundleLang) AfterResolvingDeps(context.Context) {
	ext.logf("Lifecycle: AfterResolvingDeps() called - flushing and closing log file")
	if ext.logWriter != nil {
		// Sync flushes any buffered data to disk
		_ = ext.logWriter.Sync()
		// Close the file
		_ = ext.logWriter.Close()
		ext.logWriter = nil
	}
}
