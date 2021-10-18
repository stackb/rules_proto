package protoc

import (
	"errors"
	"fmt"
	"log"
	"path"
	"sort"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
)

const (
	ResolverLangName = "protobuf"
)

var (
	errSkipImport = errors.New("self import")
	errNotFound   = errors.New("rule not found")
)

type DepsResolver func(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label)

func ResolveDepsAttr(attrName string, excludeWkt bool) DepsResolver {
	return func(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
		debug := false

		existing := r.AttrStrings(attrName)
		r.DelAttr(attrName)

		depSet := make(map[string]bool)
		for _, d := range existing {
			depSet[d] = true
		}

		for _, imp := range imports {
			if debug {
				log.Println(from, "resolving:", imp)
			}
			if excludeWkt && strings.HasPrefix(imp, "google/protobuf/") {
				continue
			}
			l, err := resolveAnyKind(c, ix, r, imp, from)
			if err == errSkipImport {
				if debug {
					log.Println(from, "skipped:", imp)
				}
				continue
			} else if err != nil {
				log.Print(err)
			} else {
				if l != label.NoLabel {
					l = l.Rel(from.Repo, from.Pkg)
					if debug {
						log.Println(from, "resolved:", imp, l)
					}
					depSet[l.String()] = true
				} else {
					if debug {
						log.Println(from, "no label", imp)
					}
				}
			}
		}

		if len(depSet) > 0 {
			deps := make([]string, 0, len(depSet))
			for dep := range depSet {
				deps = append(deps, dep)
			}
			sort.Strings(deps)
			r.SetAttr(attrName, deps)
		}
	}
}

func resolveAnyKind(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imp string, from label.Label) (label.Label, error) {
	if l, ok := resolve.FindRuleWithOverride(c, resolve.ImportSpec{Lang: r.Kind(), Imp: imp}, ResolverLangName); ok {
		// log.Println(from, "override hit:", l)
		return l, nil
	}
	if l, err := resolveWithIndex(c, ix, r.Kind(), imp, from); err == nil || err == errSkipImport {
		return l, err
	} else if err != errNotFound {
		return label.NoLabel, err
	}
	// // if debug {
	// log.Println(from, "fallback miss:", imp)
	// // }
	return label.NoLabel, nil
}

func resolveWithIndex(c *config.Config, ix *resolve.RuleIndex, kind, imp string, from label.Label) (label.Label, error) {
	matches := ix.FindRulesByImportWithConfig(c, resolve.ImportSpec{Lang: kind, Imp: imp}, ResolverLangName)
	if len(matches) == 0 {
		// log.Println(from, "no matches:", imp)
		return label.NoLabel, errNotFound
	}
	if len(matches) > 1 {
		return label.NoLabel, fmt.Errorf("multiple rules (%s and %s) may be imported with %q from %s", matches[0].Label, matches[1].Label, imp, from)
	}
	if matches[0].IsSelfImport(from) {
		// log.Println(from, "self import:", imp)
		return label.NoLabel, errSkipImport
	}
	// log.Println(from, "FindRulesByImportWithConfig: first match:", imp, matches[0].Label)
	return matches[0].Label, nil
}

// StripRel removes the rel prefix from a filename (if has matching prefix)
func StripRel(rel string, filename string) string {
	if !strings.HasPrefix(filename, rel) {
		return filename
	}
	filename = filename[len(rel):]
	return strings.TrimPrefix(filename, "/")
}

func ProtoLibraryImportSpecsForKind(kind string, libs ...ProtoLibrary) []resolve.ImportSpec {
	specs := make([]resolve.ImportSpec, 0)

	for _, lib := range libs {
		files := lib.Files()
		for _, file := range files {
			specs = append(specs, resolve.ImportSpec{Lang: kind, Imp: path.Join(file.Dir, file.Basename)})
		}
	}

	return specs
}
