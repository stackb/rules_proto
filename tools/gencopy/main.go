// gencopy is a utility program that copies bazel outputs back into the
// workspace source tree.  Ideally, you don't have any generated files committed
// to VCS, but sometimes you do.
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/google/go-cmp/cmp"
)

const (
	ModeUpdate = "update"
	ModeCheck  = "check"
)

var (
	config       = flag.String("config", "", "The JSON configuration file")
	workspaceDir = flag.String("workspace_root_directory", "", "The absolute path to the workspace source root")
)

type (
	// Config can be produced by a starlark struct.to_json() using camelCase
	// names.
	Config struct {
		// The root of the monorepo.  This comes from the environment variable
		// BUILD_WORKSPACE_DIRECTORY during a `bazel run`
		WorkspaceRootDirectory string
		// The label triggering this run
		TargetLabel string
		// The label name used for the 'update' mode
		UpdateTargetLabelName string
		// The directory name where the files were generated
		TargetPackage string
		// The list of files that were generated in the bazel output tree.  These
		// should be absolute paths.
		GeneratedFiles []string
		// The list of files that exist in the source file tree.  These are only
		// considered when the mode is 'check'
		SourceFiles []string
		// By default gencopy will perform file copy from source to destination.  If
		// mode == "check", a file difference check will be performed to
		// assert that the source and dst file contents are identical.
		Mode string
	}

	srcDst struct {
		src, dst string
	}
)

// fileExists returns whether a particular filename exists.
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// copyFile - copy bytes from one file to another
func copyFile(src, dst string) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("CopyFile: srcFile not found: %s", src)
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}

	return out.Close()
}

// readFileAsString reads the given file assumed to be text
func readFileAsString(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("could not read %s: %v", filename, err)
	}
	return string(bytes), nil
}

func usageHint(cfg *Config) string {
	return fmt.Sprintf(`You may need to regenerate the files (bazel run) using the '.%[2]s' target,
update the 'srcs = [...]' attribute to include the generated files, and then re-run the test (bazel test):

$ bazel run %[1]s.%[2]s
$ bazel test %[1]s

`, cfg.TargetLabel, cfg.UpdateTargetLabelName)
}

func check(cfg *Config, pairs []*srcDst) error {
	lenGen := len(pairs)
	lenSrc := len(cfg.SourceFiles)
	if lenSrc != lenGen {
		return fmt.Errorf(
			"check failed.  The number of source files (%d) does not match the number of generated files (%d)\n\n%s",
			lenSrc, lenGen, usageHint(cfg))
	}
	for i, pair := range pairs {
		expected, err := readFileAsString(pair.dst)
		if err != nil {
			return fmt.Errorf("check failed: %v", err)
		}
		actual, err := readFileAsString(cfg.SourceFiles[i])
		if err != nil {
			return fmt.Errorf("check failed: %v", err)
		}
		if diff := cmp.Diff(expected, actual); diff != "" {
			return fmt.Errorf("gencopy mismatch (-want +got):\n%s", diff)
		}
	}

	log.Printf("Target %s: generated files are up-to-date:", cfg.TargetLabel)
	for _, filename := range cfg.SourceFiles {
		log.Printf("  %s", filename)
	}

	return nil
}

func update(cfg *Config, pairs []*srcDst) error {
	for _, pair := range pairs {
		if err := os.MkdirAll(filepath.Base(pair.dst), os.ModePerm); err != nil {
			return fmt.Errorf("could not copy file (directory create error): %w", err)
		}
		if err := copyFile(pair.src, pair.dst); err != nil {
			return fmt.Errorf("could not copy file (%v): %w", pair, err)
		}
	}

	log.Printf("Target %s: output files copied to source tree:", cfg.TargetLabel)
	for _, pair := range pairs {
		log.Printf("  %s", pair.dst[len(cfg.WorkspaceRootDirectory)+1:])
	}

	return nil
}

func run(cfg *Config) error {
	// Prepare the src -> dst pairs
	pairs := make([]*srcDst, 0)
	for _, src := range cfg.GeneratedFiles {
		if !fileExists(src) {
			return fmt.Errorf("could not prepare (file not found): %q", src)
		}
		base := filepath.Base(src)
		dst := filepath.Join(cfg.WorkspaceRootDirectory, cfg.TargetPackage, base)
		pairs = append(pairs, &srcDst{src, dst})
	}

	switch cfg.Mode {
	case ModeCheck:
		return check(cfg, pairs)
	case ModeUpdate:
		return update(cfg, pairs)
	default:
		return fmt.Errorf("unknown run mode %q (should be one of %s, %s", cfg.Mode, ModeCheck, ModeUpdate)
	}
}

func readConfig(workspaceRootDirectory string) (*Config, error) {
	// if workspaceRootDirectory == "" {
	// 	return nil, fmt.Errorf("--workspace_root_directory is required.")
	// }
	data, err := ioutil.ReadFile(*config)
	if err != nil {
		return nil, fmt.Errorf("could not read config file %s: %w", *config, err)
	}

	cfg := &Config{
		Mode:                   "update",
		WorkspaceRootDirectory: workspaceRootDirectory,
	}

	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("could not parse config file: %w", err)
	}

	// log.Printf("%+v", cfg)

	return cfg, nil
}

func main() {
	flag.Parse()

	cfg, err := readConfig(*workspaceDir)
	if err != nil {
		log.Fatalf("gencopy: %v", err)
	}

	if err := run(cfg); err != nil {
		log.Fatalf("gencopy: %v", err)
	}
}
