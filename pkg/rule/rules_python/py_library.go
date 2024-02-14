package rules_python

import (
	"path/filepath"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

var pyLibraryKindInfo = rule.KindInfo{
	MergeableAttrs: map[string]bool{
		"srcs":       true,
		"deps":       true,
		"visibility": true,
		"imports":    true,
	},
	NonEmptyAttrs: map[string]bool{
		"srcs": true,
	},
	ResolveAttrs: map[string]bool{
		"deps": true,
	},
}

// PyLibrary implements RuleProvider for 'py_library'-derived rules.
type PyLibrary struct {
	KindName       string
	RuleNameSuffix string
	Outputs        []string
	Config         *protoc.ProtocConfiguration
	RuleConfig     *protoc.LanguageRuleConfig
	Resolver       protoc.DepsResolver
}

// Kind implements part of the ruleProvider interface.
func (s *PyLibrary) Kind() string {
	return s.KindName
}

// Name implements part of the ruleProvider interface.
func (s *PyLibrary) Name() string {
	return s.Config.Library.BaseName() + s.RuleNameSuffix
}

// Srcs computes the srcs list for the rule.
func (s *PyLibrary) Srcs() []string {
	srcs := make([]string, 0)
	for _, output := range s.Outputs {
		if strings.HasSuffix(output, ".py") {
			srcs = append(srcs, protoc.StripRel(s.Config.Rel, output))
		}
	}
	return srcs
}

// Deps computes the deps list for the rule.
func (s *PyLibrary) Deps() []string {
	return s.RuleConfig.GetDeps()
}

// Visibility provides visibility labels.
func (s *PyLibrary) Visibility() []string {
	return s.RuleConfig.GetVisibility()
}

// ImportsAttr provides the py_library.imports attribute values.
func (s *PyLibrary) ImportsAttr() (imps []string) {
	// if we have a strip_import_prefix on the proto_library, the python search
	// path should include the directory N parents above the current package,
	// where N is the number of segments needed to ascend to the prefix from
	// the dir for the current rule.
	if s.Config.Library.StripImportPrefix() == "" {
		return
	}
	prefix := s.Config.Library.StripImportPrefix()
	if !strings.HasPrefix(prefix, "/") {
		return // deal with relative-imports at another time
	}

	prefix = strings.TrimPrefix(prefix, "/")
	rel, err := filepath.Rel(prefix, s.Config.Rel)
	if err != nil {
		return // the prefix doesn't prefix the current path, shouldn't happen
	}

	parts := strings.Split(rel, "/")
	for i := 0; i < len(parts); i++ {
		parts[i] = ".."
	}
	imp := strings.Join(parts, "/")
	imps = append(imps, imp)
	return
}

// Rule implements part of the ruleProvider interface.
func (s *PyLibrary) Rule(otherGen ...*rule.Rule) *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	newRule.SetAttr("srcs", s.Srcs())

	deps := s.Deps()
	if len(deps) > 0 {
		newRule.SetAttr("deps", deps)
	}

	imports := s.ImportsAttr()
	if len(imports) > 0 {
		newRule.SetAttr("imports", imports)
	}

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	return newRule
}

// Imports implements part of the RuleProvider interface.
func (s *PyLibrary) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	if lib, ok := r.PrivateAttr(protoc.ProtoLibraryKey).(protoc.ProtoLibrary); ok {
		specs := protoc.ProtoLibraryImportSpecsForKind(r.Kind(), lib)
		specs = maybeStripImportPrefix(specs, lib.StripImportPrefix())
		return specs
	}
	return nil
}

// Resolve implements part of the RuleProvider interface.
func (s *PyLibrary) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
	s.Resolver(c, ix, r, imports, from)
}

func maybeStripImportPrefix(specs []resolve.ImportSpec, stripImportPrefix string) []resolve.ImportSpec {
	if stripImportPrefix == "" {
		return specs
	}

	prefix := strings.TrimPrefix(stripImportPrefix, "/")
	for i, spec := range specs {
		spec.Imp = strings.TrimPrefix(spec.Imp, prefix)
		spec.Imp = strings.TrimPrefix(spec.Imp, "/") // should never be absolute
		specs[i] = spec
	}

	return specs
}
