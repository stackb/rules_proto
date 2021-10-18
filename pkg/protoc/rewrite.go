package protoc

import (
	"fmt"
	// "log"
	"regexp"
	"strings"
)

// Rewrite is a replacement specification
type Rewrite struct {
	Match   *regexp.Regexp
	Replace string
}

func (r *Rewrite) ReplaceAll(src string) string {
	if r.Match == nil {
		return ""
	}
	replaced := r.Match.ReplaceAllString(src, r.Replace)
	if src == replaced {
		return ""
	}
	return replaced
}

func ParseRewrite(spec string) (*Rewrite, error) {
	parts := strings.Fields(spec)
	if len(parts) != 2 {
		return nil, fmt.Errorf("rewrite specification should be two space-separated fields [REGEXP REPLACEMENT], got %d: %v", len(parts), parts)
	}
	match, err := regexp.Compile(parts[0])
	if err != nil {
		return nil, err
	}
	return &Rewrite{match, parts[1]}, nil
}

// ResolveRewrites takes a list of rewrite rules and returns the first match of
// the given input string.
func ResolveRewrites(rewrites []Rewrite, in string) string {
	if len(rewrites) == 0 {
		return ""
	}
	for _, rw := range rewrites {
		if match := rw.ReplaceAll(in); match != "" {
			// log.Printf("SUCCESS match rewrite %q for %v", in, rw)
			return match
		}
		// log.Printf("failed to match rewrite %q for %v", in, rw)
	}
	return ""
}

// ResolveFileRewrites takes a proto File object and returns a list of matching
// rewrites of the import statements.  The list if neither sorted or
// deduplicated.
func ResolveFileRewrites(rewrites []Rewrite, file *File) []string {
	if len(rewrites) == 0 {
		return nil
	}
	resolved := make([]string, 0)
	for _, i := range file.Imports() {
		if m := ResolveRewrites(rewrites, i.Filename); m != "" {
			resolved = append(resolved, m)
		}
	}
	return resolved
}

// ResolveLibraryRewrites takes a proto_library object and returns a list
// of matching rewrites of all the the transitive import statements.
func ResolveLibraryRewrites(rewrites []Rewrite, library ProtoLibrary) []string {
	if len(rewrites) == 0 {
		return nil
	}
	resolved := make([]string, 0)
	for _, file := range library.Files() {
		resolved = append(resolved, ResolveFileRewrites(rewrites, file)...)
	}
	return DeduplicateAndSort(resolved)
}
