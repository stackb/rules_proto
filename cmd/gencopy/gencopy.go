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
	"strconv"
	"strings"

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
		// FileMode is the desired FileMode for file copy operations.  This
		// should be expressed as an octal number like '0644'.
		FileMode string
		// The set of packages we are generating for
		PackageConfigs []*PackageConfig
		// An optional file extension to append to the copied file
		Extension string
	}

	PackageConfig struct {
		// The label triggering this run
		TargetLabel string
		// The directory name where the files were generated
		TargetPackage string
		// TargetWorkspaceRoot is the value of Label.workspace_root.
		TargetWorkspaceRoot string
		// The list of files that were generated in the bazel output tree.  These
		// should be absolute paths.
		GeneratedFiles []string
		// The list of files that exist in the source file tree.  These are only
		// considered when the mode is 'check'
		SourceFiles []string
	}

	SrcDst struct {
		Src, Dst string
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
func copyFile(src, dst string, mode os.FileMode) error {
	if _, err := os.Stat(src); os.IsNotExist(err) {
		return fmt.Errorf("copyFile: src not found: %s", src)
	}

	// NOTE: for some reason the io.Copy approach was writing an empty file...
	// for now OK to copy in-memory

	data, err := ioutil.ReadFile(src)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	return ioutil.WriteFile(dst, data, mode)
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

func check(cfg *Config, pkg *PackageConfig, pairs []*SrcDst) error {
	for _, pair := range pairs {
		expected, err := readFileAsString(pair.Src)
		if err != nil {
			return fmt.Errorf("check failed while reading src %s: %v", pair.Src, err)
		}
		actual, err := readFileAsString(pair.Dst)
		if err != nil {
			return fmt.Errorf("check failed while reading dst %s: %v", pair.Dst, err)
		}

		if diff := cmp.Diff(expected, actual); diff != "" {
			return fmt.Errorf("gencopy mismatch %q vs. %q (-want +got):\n%s", pair.Src, pair.Dst, diff)
		}
	}

	fmt.Printf("Target %s: generated files are up-to-date:\n", pkg.TargetLabel)
	for _, pair := range pairs {
		fmt.Printf("  %s\n", pair.Dst)
	}

	return nil
}

func update(cfg *Config, pkg *PackageConfig, pairs []*SrcDst) error {
	for _, pair := range pairs {
		pair.Dst += cfg.Extension
	}

	mode, err := parseFileMode(cfg.FileMode)
	if err != nil {
		return fmt.Errorf("update: %v", err)
	}
	for _, pair := range pairs {
		if err := os.MkdirAll(filepath.Base(pair.Dst), os.ModePerm); err != nil {
			return fmt.Errorf("could not copy file (directory create error): %w", err)
		}
		if err := copyFile(pair.Src, pair.Dst, mode); err != nil {
			return fmt.Errorf("could not copy file pair (%+v): %w", pair, err)
		}
	}

	fmt.Printf("Target %s: output files copied to source tree:\n", pkg.TargetLabel)
	for _, pair := range pairs {
		fmt.Printf("  %s\n", pair.Dst[len(cfg.WorkspaceRootDirectory)+1:])
	}

	return nil
}

func makePkgSrcDstPairs(cfg *Config, pkg *PackageConfig) []*SrcDst {
	// Prepare the src -> dst pairs
	pairs := make([]*SrcDst, len(pkg.GeneratedFiles))

	// we are copying/comparing generated files to their source file
	// equivalents.  So here 'src' is the generated file and 'dst' is the
	// source file target. So yeah.
	for i, src := range pkg.GeneratedFiles {
		pairs[i] = makePkgSrcDstPair(cfg, pkg, src, pkg.SourceFiles[i])
	}

	return pairs
}

func makePkgSrcDstPair(cfg *Config, pkg *PackageConfig, src, dst string) *SrcDst {
	if pkg.TargetWorkspaceRoot != "" {
		src = filepath.Join("external", strings.TrimPrefix(src, ".."))
		dst = filepath.Join(pkg.TargetWorkspaceRoot, dst)
	}
	dst = filepath.Join(cfg.WorkspaceRootDirectory, dst)
	return &SrcDst{src, dst}
}

func runPkg(cfg *Config, pkg *PackageConfig) (err error) {
	pairs := makePkgSrcDstPairs(cfg, pkg)

	for _, pair := range pairs {
		if !fileExists(pair.Src) {
			return fmt.Errorf("could not prepare (generated file not found): %q", pair.Src)
		}
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

	if cfg.FileMode == "" {
		cfg.FileMode = "0644"
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

func parseFileMode(s string) (os.FileMode, error) {
	value, err := strconv.ParseUint(s, 0, 32)
	if err != nil {
		return 0, err
	}
	return os.FileMode(value), nil
}
