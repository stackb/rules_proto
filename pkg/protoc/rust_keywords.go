package protoc

import (
	"path"
	"strings"
)

var rustKeywords = map[string]bool{
	// Strict keywords
	"as":       true,
	"break":    true,
	"const":    true,
	"continue": true,
	"crate":    true,
	"else":     true,
	"enum":     true,
	"extern":   true,
	"false":    true,
	"fn":       true,
	"for":      true,
	"if":       true,
	"impl":     true,
	"in":       true,
	"let":      true,
	"loop":     true,
	"match":    true,
	"mod":      true,
	"move":     true,
	"mut":      true,
	"pub":      true,
	"ref":      true,
	"return":   true,
	"self":     true,
	"Self":     true,
	"static":   true,
	"struct":   true,
	"super":    true,
	"trait":    true,
	"true":     true,
	"type":     true,
	"unsafe":   true,
	"use":      true,
	"where":    true,
	"while":    true,
	// Async keywords (edition 2018+)
	"async": true,
	"await": true,
	"dyn":   true,
	// Reserved keywords
	"abstract": true,
	"become":   true,
	"box":      true,
	"do":       true,
	"final":    true,
	"macro":    true,
	"override": true,
	"priv":     true,
	"try":      true,
	"typeof":   true,
	"unsized":  true,
	"virtual":  true,
	"yield":    true,
}

// RustProtocOutputDir returns the directory that protoc-gen-prost (and its
// siblings) will write outputs into for a given proto package, with Rust
// keyword segments escaped via the r# prefix.
//
// Examples:
//   - "google.type"                    → "google/r#type"
//   - "trumid.common.auth"             → "trumid/common/auth"
//   - ""                                → ""
//
// Note this is independent of the bazel package location — protoc-gen-prost
// always derives the output path from the proto package name unless given
// flat_output_dir=true. Use the result to compare against pc.Rel and decide
// whether output_mappings are needed.
func RustProtocOutputDir(pkg string) string {
	if pkg == "" {
		return ""
	}
	segments := strings.Split(pkg, ".")
	for i, seg := range segments {
		if rustKeywords[seg] {
			segments[i] = "r#" + seg
		}
	}
	return strings.Join(segments, "/")
}

// RustKeywordEscapeMappings computes output mappings needed when
// protoc-gen-prost escapes Rust keywords with the r# prefix in directory paths.
//
// For example, proto package "google.type" causes prost to write files to
// "google/r#type/" instead of "google/type/". This function returns a mapping
// from each declared output filename to the actual prost output path.
//
// Returns an empty map if no package segments are Rust keywords. Callers who
// also need to handle the more general case of proto-package path differing
// from the bazel package path should use RustProtocOutputDir directly.
func RustKeywordEscapeMappings(pkg string, outputs []string) map[string]string {
	if pkg == "" || len(outputs) == 0 {
		return nil
	}

	// Check if any segment is a Rust keyword.
	needsEscape := false
	for _, seg := range strings.Split(pkg, ".") {
		if rustKeywords[seg] {
			needsEscape = true
			break
		}
	}
	if !needsEscape {
		return nil
	}

	escapedDir := RustProtocOutputDir(pkg)
	mappings := make(map[string]string, len(outputs))
	for _, output := range outputs {
		base := path.Base(output)
		mappings[base] = path.Join(escapedDir, base)
	}
	return mappings
}

// RustCrateName returns the canonical Rust crate name for a proto package.
// The proto package's dots are replaced with underscores; the crate name IS
// the namespace, so no language suffix is appended (e.g.
// "trumid.common.utils.state.snapshot.proto" →
// "trumid_common_utils_state_snapshot_proto"). The proto_rust_library macro
// is expected to expose all generated types at the crate root, matching this
// flat convention. Returns the empty string for an empty input.
func RustCrateName(protoPackage string) string {
	if protoPackage == "" {
		return ""
	}
	return strings.ReplaceAll(protoPackage, ".", "_")
}
