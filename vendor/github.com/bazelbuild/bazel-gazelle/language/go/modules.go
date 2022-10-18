/* Copyright 2018 The Bazel Authors. All rights reserved.

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

package golang

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

func importReposFromModules(args language.ImportReposArgs) language.ImportReposResult {
	dir := filepath.Dir(args.Path)

	// List all modules except for the main module, including implicit indirect
	// dependencies.
	type module struct {
		Path, Version, Sum, Error string
		Main               bool
		Replace            *struct {
			Path, Version string
		}
	}
	// path@version can be used as a unique identifier for looking up sums
	pathToModule := map[string]*module{}
	data, err := goListModules(dir)
	if err != nil {
		return language.ImportReposResult{Error: err}
	}
	dec := json.NewDecoder(bytes.NewReader(data))
	for dec.More() {
		mod := new(module)
		if err := dec.Decode(mod); err != nil {
			return language.ImportReposResult{Error: err}
		}
		if mod.Main {
			continue
		}
		if mod.Replace != nil {
			if filepath.IsAbs(mod.Replace.Path) || build.IsLocalImport(mod.Replace.Path) {
				log.Printf("go_repository does not support file path replacements for %s -> %s", mod.Path,
					mod.Replace.Path)
				continue
			}
			pathToModule[mod.Replace.Path+"@"+mod.Replace.Version] = mod
		} else {
			pathToModule[mod.Path+"@"+mod.Version] = mod
		}
	}

	// Load sums from go.sum. Ideally, they're all there.
	goSumPath := filepath.Join(filepath.Dir(args.Path), "go.sum")
	data, _ = ioutil.ReadFile(goSumPath)
	lines := bytes.Split(data, []byte("\n"))
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		fields := bytes.Fields(line)
		if len(fields) != 3 {
			continue
		}
		path, version, sum := string(fields[0]), string(fields[1]), string(fields[2])
		if strings.HasSuffix(version, "/go.mod") {
			continue
		}
		if mod, ok := pathToModule[path+"@"+version]; ok {
			mod.Sum = sum
		}
	}

	// If sums are missing, run 'go mod download' to get them.
	// This must be done in a temporary directory because 'go mod download'
	// may modify go.mod and go.sum. It does not support -mod=readonly.
	var missingSumArgs []string
	for pathVer, mod := range pathToModule {
		if mod.Sum == "" {
			missingSumArgs = append(missingSumArgs, pathVer)
		}
	}

	if len(missingSumArgs) > 0 {
		tmpDir, err := ioutil.TempDir("", "")
		if err != nil {
			return language.ImportReposResult{Error: fmt.Errorf("finding module sums: %v", err)}
		}
		defer os.RemoveAll(tmpDir)
		data, err := goModDownload(tmpDir, missingSumArgs)
		dec = json.NewDecoder(bytes.NewReader(data))
		if err != nil {
			// Best-effort try to adorn specific error details from the JSON output.
			for dec.More() {
				var dl module
				if err := dec.Decode(&dl); err != nil {
					// If we couldn't parse a possible error description, just ignore this part of the output.
					continue
				}
				if dl.Error != "" {
					err = fmt.Errorf("%v\nError downloading %v: %v", err, dl.Path, dl.Error)
				}
			}

			return language.ImportReposResult{Error: err}
		}
		for dec.More() {
			var dl module
			if err := dec.Decode(&dl); err != nil {
				return language.ImportReposResult{Error: err}
			}
			if mod, ok := pathToModule[dl.Path+"@"+dl.Version]; ok {
				mod.Sum = dl.Sum
			}
		}
	}

	// Translate to repository rules.
	gen := make([]*rule.Rule, 0, len(pathToModule))
	for pathVer, mod := range pathToModule {
		if mod.Sum == "" {
			log.Printf("could not determine sum for module %s", pathVer)
			continue
		}
		r := rule.NewRule("go_repository", label.ImportPathToBazelRepoName(mod.Path))
		r.SetAttr("importpath", mod.Path)
		r.SetAttr("sum", mod.Sum)
		if mod.Replace == nil {
			r.SetAttr("version", mod.Version)
		} else {
			r.SetAttr("replace", mod.Replace.Path)
			r.SetAttr("version", mod.Replace.Version)
		}
		gen = append(gen, r)
	}
	sort.Slice(gen, func(i, j int) bool {
		return gen[i].Name() < gen[j].Name()
	})
	return language.ImportReposResult{Gen: gen}
}

// goListModules invokes "go list" in a directory containing a go.mod file.
var goListModules = func(dir string) ([]byte, error) {
	return runGoCommandForOutput(dir, "list", "-mod=readonly", "-e", "-m", "-json", "all")
}

// goModDownload invokes "go mod download" in a directory containing a
// go.mod file.
var goModDownload = func(dir string, args []string) ([]byte, error) {
	dlArgs := []string{"mod", "download", "-json"}
	dlArgs = append(dlArgs, args...)
	return runGoCommandForOutput(dir, dlArgs...)
}

// findGoTool attempts to locate the go executable. If GOROOT is set, we'll
// prefer the one in there; otherwise, we'll rely on PATH. If the wrapper
// script generated by the gazelle rule is invoked by Bazel, it will set
// GOROOT to the configured SDK. We don't want to rely on the host SDK in
// that situation.
func findGoTool() string {
	path := "go" // rely on PATH by default
	if goroot, ok := os.LookupEnv("GOROOT"); ok {
		path = filepath.Join(goroot, "bin", "go")
	}
	if runtime.GOOS == "windows" {
		path += ".exe"
	}
	return path
}

func runGoCommandForOutput(dir string, args ...string) ([]byte, error) {
	goTool := findGoTool()
	env := os.Environ()
	env = append(env, "GO111MODULE=on")
	if os.Getenv("GOCACHE") == "" && os.Getenv("HOME") == "" {
		gocache, err := ioutil.TempDir("", "")
		if err != nil {
			return nil, err
		}
		env = append(env, "GOCACHE="+gocache)
		defer os.RemoveAll(gocache)
	}
	if os.Getenv("GOPATH") == "" && os.Getenv("HOME") == "" {
		gopath, err := ioutil.TempDir("", "")
		if err != nil {
			return nil, err
		}
		env = append(env, "GOPATH="+gopath)
		defer os.RemoveAll(gopath)
	}
	cmd := exec.Command(goTool, args...)
	stderr := &bytes.Buffer{}
	cmd.Stderr = stderr
	cmd.Dir = dir
	cmd.Env = env
	out, err := cmd.Output()
	if err != nil {
		var errStr string
		var xerr *exec.ExitError
		if errors.As(err, &xerr) {
			errStr = strings.TrimSpace(stderr.String())
		} else {
			errStr = err.Error()
		}
		return out, fmt.Errorf("running '%s %s': %s", cmd.Path, strings.Join(cmd.Args, " "), errStr)
	}
	return out, nil
}
