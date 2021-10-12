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

func ResolveDepsAttr(attrName string) DepsResolver {
	return ResolveDepsAttrDebug(attrName, false)
}

func ResolveDepsAttrDebug(attrName string, debug bool) DepsResolver {
	return func(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
		existing := r.AttrStrings(attrName)
		r.DelAttr(attrName)
		debug = true

		depSet := make(map[string]bool)
		for _, d := range existing {
			depSet[d] = true
		}

		for _, imp := range imports {
			if debug {
				log.Println(from, "resolving:", imp)
			}
			if strings.HasPrefix(imp, "google/protobuf/") {
				continue
			}
			l, err := resolveAnyKind(c, ix, r, imp, from, debug)
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
					// log.Println(from, "no label:", imp)
				}
			}

			// result := cfg.Resolve(r.Kind(), "srcs", imp)
			// if len(result) > 0 {
			// 	first := result[0]
			// 	deps = append(deps, first.Label.String())
			// }
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

func resolveAnyKind(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imp string, from label.Label, debug bool) (label.Label, error) {
	if l, ok := resolve.FindRuleWithOverride(c, resolve.ImportSpec{Lang: r.Kind(), Imp: imp}, ResolverLangName); ok {
		return l, nil
	}
	// if match := GlobalResolver().Resolve(r.Kind(), "imports", imp); len(match) > 0 {
	// 	if l.Equal(from) {
	// 		return label.NoLabel, errSkipImport
	// 	} else {
	// 		return l, nil
	// 	}
	// }
	if l, err := resolveWithIndex(c, ix, r.Kind(), imp, from); err == nil || err == errSkipImport {
		return l, err
	} else if err != errNotFound {
		if debug {
			log.Println(from, "error:", imp, err)
		}
		return label.NoLabel, err
	}
	if debug {
		log.Println(from, "fallback miss:", imp)
	}
	return label.NoLabel, nil
}

func resolveWithIndex(c *config.Config, ix *resolve.RuleIndex, kind, imp string, from label.Label) (label.Label, error) {
	matches := ix.FindRulesByImportWithConfig(c, resolve.ImportSpec{Lang: kind, Imp: imp}, ResolverLangName)
	if len(matches) == 0 {
		return label.NoLabel, errNotFound
	}
	if len(matches) > 1 {
		return label.NoLabel, fmt.Errorf("multiple rules (%s and %s) may be imported with %q from %s", matches[0].Label, matches[1].Label, imp, from)
	}
	if matches[0].IsSelfImport(from) {
		return label.NoLabel, errSkipImport
	}
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
