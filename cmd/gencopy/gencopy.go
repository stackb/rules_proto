// gencopy is a utility program that copies bazel outputs back into the
// workspace source tree.  Ideally, you don't have any generated files committed
// to VCS, but sometimes you do.
//
package main

import (
	"encoding/json"
	"flag"
	"fmt"
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
		// The label name used for the 'update' mode
		UpdateTargetLabelName string
		// By default gencopy will perform file copy from source to destination.  If
		// mode == "check", a file difference check will be performed to
		// assert that the source and dst file contents are identical.
		Mode string
		// The set of packages we are generating for
		PackageConfigs []*PackageConfig
	}

	PackageConfig struct {
		// The label triggering this run
		TargetLabel string
		// The directory name where the files were generated
		TargetPackage string
		// The list of files that were generated in the bazel output tree.  These
		// should be absolute paths.
		GeneratedFiles []string
		// The list of files that exist in the source file tree.  These are only
		// considered when the mode is 'check'
		SourceFiles []string
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
		return fmt.Errorf("copyFile: src not found: %s", src)
	}

	// NOTE: for some reason the io.Copy approach was writing an empty file...
	// for now OK to copy in-memory

	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(dst, data, os.ModePerm)

}

// readFileAsString reads the given file assumed to be text
func readFileAsString(filename string) (string, error) {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return "", fmt.Errorf("could not read %s: %v", filename, err)
	}
	return string(bytes), nil
}

func usageHint(cfg *Config, pkg *PackageConfig) string {
	return fmt.Sprintf(`You may need to regenerate the files (bazel run) using the '.%[2]s' target,
update the 'srcs = [...]' attribute to include the generated files and re-run the test:

$ bazel run %[1]s.%[2]s
$ bazel test %[1]s

`, pkg.TargetLabel, cfg.UpdateTargetLabelName)
}

func check(cfg *Config, pkg *PackageConfig, pairs []*srcDst) error {
	for _, pair := range pairs {
		expected, err := readFileAsString(pair.src)
		if err != nil {
			return fmt.Errorf("check failed while reading src %s: %v", pair.src, err)
		}
		actual, err := readFileAsString(pair.dst)
		if err != nil {
			return fmt.Errorf("check failed while reading dst %s: %v", pair.dst, err)
		}
		if diff := cmp.Diff(expected, actual); diff != "" {
			return fmt.Errorf("gencopy mismatch %q vs. %q (-want +got):\n%s", pair.src, pair.dst, diff)
		}
	}

	fmt.Printf("Target %s: generated files are up-to-date:\n", pkg.TargetLabel)
	for _, pair := range pairs {
		fmt.Printf("  %s\n", pair.dst)
	}

	return nil
}

func update(cfg *Config, pkg *PackageConfig, pairs []*srcDst) error {
	for _, pair := range pairs {
		if err := os.MkdirAll(filepath.Base(pair.dst), os.ModePerm); err != nil {
			return fmt.Errorf("could not copy file (directory create error): %w", err)
		}
		if err := copyFile(pair.src, pair.dst); err != nil {
			return fmt.Errorf("could not copy file pair (%+v): %w", pair, err)
		}
	}

	fmt.Printf("Target %s: output files copied to source tree:\n", pkg.TargetLabel)
	for _, pair := range pairs {
		fmt.Printf("  %s\n", pair.dst[len(cfg.WorkspaceRootDirectory)+1:])
	}

	return nil
}

func runPkg(cfg *Config, pkg *PackageConfig) (err error) {
	// Prepare the src -> dst pairs
	pairs := make([]*srcDst, len(pkg.GeneratedFiles))

	// we are copying/comparing generated files to their source file
	// equivalents.  So here 'src' is the generated file and 'dst' is the
	// source file target. So yeah.
	for i, src := range pkg.GeneratedFiles {
		if !fileExists(src) {
			return fmt.Errorf("could not prepare (generated file not found): %q", src)
		}
		dst := filepath.Join(cfg.WorkspaceRootDirectory, pkg.SourceFiles[i])
		pair := &srcDst{src, dst}
		pairs[i] = pair
		// log.Printf("pair: %+v", pair)
	}

	switch cfg.Mode {
	case ModeCheck:
		err = check(cfg, pkg, pairs)
	case ModeUpdate:
		err = update(cfg, pkg, pairs)
	default:
		err = fmt.Errorf("unknown run mode %q (should be one of %s, %s", cfg.Mode, ModeCheck, ModeUpdate)
	}
	if err != nil {
		return err
	}

	return
}

func run(cfg *Config) error {
	for _, pkg := range cfg.PackageConfigs {
		if err := runPkg(cfg, pkg); err != nil {
			return err
		}
	}
	return nil
}

func readConfig(workspaceRootDirectory string) (*Config, error) {
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

	if cfg.Mode == ModeUpdate && cfg.WorkspaceRootDirectory == "" {
		log.Fatalln("workspace directory is mandatory in update mode")
	}

	if err := run(cfg); err != nil {
		log.Fatalf("gencopy: %v", err)
	}
}
