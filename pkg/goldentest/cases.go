package goldentest

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

import (
	"bufio"
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bazelbuild/bazel-gazelle/testtools"
	"github.com/bazelbuild/rules_go/go/tools/bazel"
)

var doCleanup = true

// GoldenTests is a helper for running glob(["testdata/**"]) style test setups.
type GoldenTests struct {
	extensionDir string
	testDataPath string
	extraArgs    []string
	dataFiles    []bazel.RunfileEntry
	onlyTests    []string
}

type GoldenTestOption func(*GoldenTests)

func WithExtraArgs(args ...string) GoldenTestOption {
	return func(g *GoldenTests) {
		g.extraArgs = args
	}
}

func WithDataFiles(files ...bazel.RunfileEntry) GoldenTestOption {
	return func(g *GoldenTests) {
		g.dataFiles = files
	}
}

func WithOnlyTests(tests ...string) GoldenTestOption {
	return func(g *GoldenTests) {
		g.onlyTests = tests
	}
}

// FromDir construct a GoldenTests tester that searches the given directory.
func FromDir(extensionDir string, options ...GoldenTestOption) *GoldenTests {
	test := &GoldenTests{
		extensionDir: extensionDir,
		testDataPath: path.Join(extensionDir, "testdata") + "/",
	}
	for _, opt := range options {
		opt(test)
	}
	return test
}

func (g *GoldenTests) Run(t *testing.T, gazelleName string) {
	t.Log("Run", g.extensionDir)
	// listFiles(".")

	gazellePath, ok := bazel.FindBinary(g.extensionDir, gazelleName)
	if !ok {
		t.Fatalf("could not find gazelle: %q in %s", gazelleName, g.extensionDir)
	}
	// t.Log("Found gazelle binary:", gazellePath)

	tests := map[string][]bazel.RunfileEntry{}

	files, err := bazel.ListRunfiles()
	if err != nil {
		t.Fatalf("bazel.ListRunfiles() error: %v", err)
	}

	for _, f := range files {
		if strings.HasPrefix(f.ShortPath, g.testDataPath) {
			relativePath := strings.TrimPrefix(f.ShortPath, g.testDataPath)
			parts := strings.SplitN(relativePath, "/", 2)
			if len(parts) < 2 {
				// This file is not a part of a testcase since it must be in a dir that
				// is the test case and then have a path inside of that.
				t.Logf("excluding file %s (not part of testcase dir)", relativePath)
				continue
			}

			tests[parts[0]] = append(tests[parts[0]], f)
		}
	}
	if len(tests) == 0 {
		t.Fatal("no tests found")
	}

	for testName, files := range tests {
		shouldTest := true
		if len(g.onlyTests) > 0 {
			shouldTest = false
			for _, name := range g.onlyTests {
				if name == testName {
					shouldTest = true
					break
				}
			}
		}
		if shouldTest {
			g.testPath(t, gazellePath, testName, files)
		} else {
			log.Println("skipped test:", testName)
		}
	}
}

func (g *GoldenTests) testPath(t *testing.T, gazellePath, name string, files []bazel.RunfileEntry) {
	t.Run(name, func(t *testing.T) {
		var inputs []testtools.FileSpec
		var goldens []testtools.FileSpec
		extraArgs := g.extraArgs

		for _, f := range files {
			path := f.Path
			trim := g.testDataPath + name + "/"
			shortPath := strings.TrimPrefix(f.ShortPath, trim)
			info, err := os.Stat(path)
			if err != nil {
				t.Fatalf("os.Stat(%q) error: %v", path, err)
			}

			// Skip dirs.
			if info.IsDir() {
				continue
			}

			content, err := ioutil.ReadFile(path)
			if err != nil {
				t.Errorf("ioutil.ReadFile(%q) error: %v", path, err)
			}

			if shortPath == ".gazelle.args" {
				extraArgs = append(extraArgs, parseArgsFile(bytes.NewReader(content))...)
				continue
			}

			// Now trim the common prefix off.
			if strings.HasSuffix(shortPath, ".in") {
				inputs = append(inputs, testtools.FileSpec{
					Path:    strings.TrimSuffix(shortPath, ".in"),
					Content: string(content),
				})
			} else if strings.HasSuffix(shortPath, ".out") {
				goldens = append(goldens, testtools.FileSpec{
					Path:    strings.TrimSuffix(shortPath, ".out"),
					Content: string(content),
				})
			} else {
				inputs = append(inputs, testtools.FileSpec{
					Path:    shortPath,
					Content: string(content),
				})
				goldens = append(goldens, testtools.FileSpec{
					Path:    shortPath,
					Content: string(content),
				})
			}
		}

		dir, cleanup := testtools.CreateFiles(t, inputs)
		if doCleanup {
			defer cleanup()
		}

		for _, f := range g.dataFiles {
			newName := filepath.Join(dir, f.ShortPath)
			newDir := filepath.Dir(newName)
			if err := os.MkdirAll(newDir, os.ModePerm); err != nil {
				t.Fatal("data file symlink setup error:", f, err)
			}
			if err := os.Symlink(f.Path, newName); err != nil {
				t.Fatal("data file symlink setup error:", f, err)
			}
		}

		t.Log("running test dir:", dir)
		args := append([]string{"-build_file_name=BUILD"}, extraArgs...)
		cmd := exec.Command(gazellePath, args...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Dir = dir
		if err := cmd.Run(); err != nil {
			t.Fatal("gazelle command failed!", err)
		}

		t.Log("checking files:", dir)

		testtools.CheckFiles(t, dir, goldens)

		if t.Failed() {
			filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return err
				}
				t.Log("file:", path)
				return nil
			})
		}
	})
}

// listFiles - convenience debugging function to log the files under a given dir
func listFiles(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("%v\n", err)
			return err
		}
		if info.Mode()&os.ModeSymlink > 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			log.Printf("%s -> %s", path, link)
			return nil
		}

		log.Println(path)
		return nil
	})
}

func parseArgsFile(in io.Reader) []string {
	args := make([]string, 0)
	scanner := bufio.NewScanner(in)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}
		args = append(args, line)
	}
	return args
}
