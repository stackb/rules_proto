// preserve_packages walks a fetched external repository and rewrites the
// upstream files that fetch_repo -clean would otherwise delete. BUILD and
// BUILD.bazel are renamed to BUILD.package / BUILD.bazel.package so the
// starlarkrepository gazelle extension can capture them as starlark_package
// rules; everything else fetch_repo -clean removes (MODULE.bazel,
// WORKSPACE, …) is deleted outright.
//
// Intended to be invoked from rules/proto/proto_repository.bzl when
// build_file_generation = "preserve".
package main

import (
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

var renameMap = map[string]string{
	"BUILD":       "BUILD.package",
	"BUILD.bazel": "BUILD.bazel.package",
}

var deleteSet = map[string]bool{
	"MODULE.bazel":      true,
	"MODULE.bazel.lock": true,
	"WORKSPACE":         true,
	"WORKSPACE.bazel":   true,
	"WORKSPACE.bzlmod":  true,
}

func main() {
	root := flag.String("root", "", "repo root to walk (required)")
	flag.Parse()
	if *root == "" {
		log.Fatal("preserve_packages: -root is required")
	}

	if err := run(*root); err != nil {
		log.Fatalf("preserve_packages: %v", err)
	}
}

func run(root string) error {
	return filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		name := info.Name()
		if dst, ok := renameMap[name]; ok {
			target := filepath.Join(filepath.Dir(path), dst)
			if err := os.Rename(path, target); err != nil {
				return fmt.Errorf("rename %s -> %s: %w", path, target, err)
			}
			return nil
		}
		if deleteSet[name] {
			if err := os.Remove(path); err != nil {
				return fmt.Errorf("remove %s: %w", path, err)
			}
		}
		return nil
	})
}
