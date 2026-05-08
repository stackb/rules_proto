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

// RustKeywordEscapeMappings computes output mappings needed when
// protoc-gen-prost escapes Rust keywords with the r# prefix in directory paths.
//
// For example, proto package "google.type" causes prost to write files to
// "google/r#type/" instead of "google/type/". This function returns a mapping
// from each declared output filename to the actual prost output path.
//
// Returns an empty map if no package segments are Rust keywords.
func RustKeywordEscapeMappings(pkg string, outputs []string) map[string]string {
	if pkg == "" || len(outputs) == 0 {
		return nil
	}

	segments := strings.Split(pkg, ".")

	// Check if any segment is a Rust keyword.
	needsEscape := false
	for _, seg := range segments {
		if rustKeywords[seg] {
			needsEscape = true
			break
		}
	}
	if !needsEscape {
		return nil
	}

	// Build the escaped directory path.
	escaped := make([]string, len(segments))
	for i, seg := range segments {
		if rustKeywords[seg] {
			escaped[i] = "r#" + seg
		} else {
			escaped[i] = seg
		}
	}
	escapedDir := strings.Join(escaped, "/")

	mappings := make(map[string]string, len(outputs))
	for _, output := range outputs {
		base := path.Base(output)
		mappings[base] = path.Join(escapedDir, base)
	}
	return mappings
}
