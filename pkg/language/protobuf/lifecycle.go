package protobuf

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Before implements part of the language.LifecycleManager interface.
func (pl *protobufLang) Before(context.Context) {
}

// DoneGeneratingRules implements part of the language.LifecycleManager interface.
//
// Performs three cross-package syncs that need every GenerateRules call to
// have completed first:
//
//  1. Root Cargo.toml [workspace] members list — the lines between the
//     `# gazelle:proto_rust_members start/end` markers are replaced with
//     one entry per package-level proto_rust_library.
//
//  2. Root BUILD.bazel proto_compile_assets aggregator deps — the lines
//     between the `# gazelle:vendor_proto_sources_deps start/end` markers
//     are replaced with one entry per generated proto_compiled_sources rule
//     and one entry per proto_rust_library's underlying _lib target.
//
//  3. Root Cargo.toml [workspace] exclude list — the lines between the
//     `# gazelle:proto_rust_excludes start/end` markers are replaced with
//     one entry per Bazel package containing standalone per-file crates.
//
// All syncs are no-ops when the corresponding markers are absent (or the
// target file does not exist).
func (pl *protobufLang) DoneGeneratingRules() {
	if pl.repoRoot == "" {
		return
	}
	if err := updateRootCargoMembers(pl.repoRoot, pl.protoRustLibraryPackages); err != nil {
		log.Printf("warning: could not update root Cargo.toml proto_rust_members: %v", err)
	}
	if err := updateRootCargoExcludes(pl.repoRoot, pl.protoRustPerFilePackageDirs); err != nil {
		log.Printf("warning: could not update root Cargo.toml proto_rust_excludes: %v", err)
	}
	if err := updateRootVendorAssetsDeps(pl.repoRoot, pl.vendorAssetLabels); err != nil {
		log.Printf("warning: could not update root BUILD.bazel vendor_proto_sources_deps: %v", err)
	}
}

// AfterResolvingDeps implements part of the language.LifecycleManager interface.
func (pl *protobufLang) AfterResolvingDeps(context.Context) {
}

const (
	cargoMembersStartMarker     = "# gazelle:proto_rust_members start"
	cargoMembersEndMarker       = "# gazelle:proto_rust_members end"
	cargoExcludesStartMarker    = "# gazelle:proto_rust_excludes start"
	cargoExcludesEndMarker      = "# gazelle:proto_rust_excludes end"
	vendorAssetsDepsStartMarker = "# gazelle:vendor_proto_sources_deps start"
	vendorAssetsDepsEndMarker   = "# gazelle:vendor_proto_sources_deps end"
)

// updateRootCargoMembers rewrites the gazelle:proto_rust_members marker
// section in the root Cargo.toml with a sorted, deduplicated list of
// `"<pkg>",` entries. No-op if the file is missing or the markers are
// absent.
func updateRootCargoMembers(repoRoot string, packages []string) error {
	return rewriteMarkerSection(
		filepath.Join(repoRoot, "Cargo.toml"),
		cargoMembersStartMarker,
		cargoMembersEndMarker,
		packages,
		"[workspace] members list",
	)
}

// updateRootCargoExcludes rewrites the gazelle:proto_rust_excludes marker
// section in the root Cargo.toml with directories containing standalone
// per-file crates.
func updateRootCargoExcludes(repoRoot string, packages []string) error {
	return rewriteMarkerSection(
		filepath.Join(repoRoot, "Cargo.toml"),
		cargoExcludesStartMarker,
		cargoExcludesEndMarker,
		packages,
		"[workspace] exclude list",
	)
}

// updateRootVendorAssetsDeps rewrites the gazelle:vendor_proto_sources_deps
// marker section in the root BUILD.bazel with a sorted, deduplicated list
// of `"<label>",` entries. No-op if the file is missing or the markers are
// absent.
func updateRootVendorAssetsDeps(repoRoot string, labels []string) error {
	return rewriteMarkerSection(
		filepath.Join(repoRoot, "BUILD.bazel"),
		vendorAssetsDepsStartMarker,
		vendorAssetsDepsEndMarker,
		labels,
		"vendor_proto_sources deps list",
	)
}

// rewriteMarkerSection replaces every line between startMarker and
// endMarker in path with a sorted, deduplicated quoted-comma list of
// entries. Marker lines themselves are preserved; indentation is taken
// from the start marker. The file is not written if the resulting content
// is identical (avoids spurious mtime bumps). When entries is non-empty
// but markers are absent, a warning is logged once with the supplied
// description so the maintainer knows where to add them.
func rewriteMarkerSection(path, startMarker, endMarker string, entries []string, description string) error {
	src, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read %s: %w", path, err)
	}

	seen := make(map[string]bool, len(entries))
	uniq := make([]string, 0, len(entries))
	for _, e := range entries {
		if seen[e] {
			continue
		}
		seen[e] = true
		uniq = append(uniq, e)
	}
	sort.Strings(uniq)

	lines := strings.Split(string(src), "\n")
	startIdx, endIdx := -1, -1
	var indent string
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if startIdx < 0 && trimmed == startMarker {
			startIdx = i
			indent = line[:len(line)-len(strings.TrimLeft(line, " \t"))]
			continue
		}
		if startIdx >= 0 && trimmed == endMarker {
			endIdx = i
			break
		}
	}
	if startIdx < 0 || endIdx < 0 {
		if len(uniq) > 0 {
			log.Printf("warning: %s has %d entries to enroll in %s but lacks the %s / %s markers — add them to enable auto-update", path, len(uniq), description, startMarker, endMarker)
		}
		return nil
	}

	newSection := make([]string, 0, len(uniq)+2)
	newSection = append(newSection, lines[startIdx])
	for _, e := range uniq {
		newSection = append(newSection, fmt.Sprintf("%s\"%s\",", indent, e))
	}
	newSection = append(newSection, lines[endIdx])

	out := append([]string{}, lines[:startIdx]...)
	out = append(out, newSection...)
	out = append(out, lines[endIdx+1:]...)

	updated := strings.Join(out, "\n")
	if updated == string(src) {
		return nil
	}
	if err := os.WriteFile(path, []byte(updated), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}
	return nil
}
