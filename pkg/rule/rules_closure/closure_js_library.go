package rules_closure

import (
	"container/list"
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const transitiveDepsKey = "_transitive_proto_library_deps"

var closureJsLibraryKindInfo = rule.KindInfo{
	MergeableAttrs: map[string]bool{
		"srcs":                 true,
		"internal_descriptors": true,
		"exports":              true,
		"visibility":           true,
		"suppress":             true,
		"lenient":              true,
	},
	ResolveAttrs: map[string]bool{"deps": true},
}

// ClosureJsLibrary implements RuleProvider for 'py_library'-derived rules.
type ClosureJsLibrary struct {
	KindName       string
	RuleNameSuffix string
	Outputs        []string
	Config         *protoc.ProtocConfiguration
	RuleConfig     *protoc.LanguageRuleConfig
	Resolver       protoc.DepsResolver
}

// Kind implements part of the ruleProvider interface.
func (s *ClosureJsLibrary) Kind() string {
	return s.KindName
}

// Name implements part of the ruleProvider interface.
func (s *ClosureJsLibrary) Name() string {
	return s.Config.Library.BaseName() + s.RuleNameSuffix
}

// Srcs computes the srcs list for the rule.
func (s *ClosureJsLibrary) Srcs() []string {
	srcs := make([]string, 0)
	for _, output := range s.Outputs {
		if strings.HasSuffix(output, "_closure.js") {
			srcs = append(srcs, protoc.StripRel(s.Config.Rel, output))
		}
	}
	return srcs
}

// Deps computes the deps list for the rule.
func (s *ClosureJsLibrary) Deps() []string {
	return s.RuleConfig.GetDeps()
}

// Visibility provides visibility labels.
func (s *ClosureJsLibrary) Visibility() []string {
	return s.RuleConfig.GetVisibility()
}

// Rule implements part of the ruleProvider interface.
func (s *ClosureJsLibrary) Rule(otherGen ...*rule.Rule) *rule.Rule {
	newRule := rule.NewRule(s.Kind(), s.Name())

	newRule.SetAttr("srcs", s.Srcs())

	visibility := s.Visibility()
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}

	newRule.SetAttr("suppress", []string{
		"JSC_IMPLICITLY_NONNULL_JSDOC",
		"JSC_UNUSED_LOCAL_ASSIGNMENT",
	})

	return newRule
}

// Imports implements part of the RuleProvider interface.
func (s *ClosureJsLibrary) Imports(c *config.Config, r *rule.Rule, file *rule.File) []resolve.ImportSpec {
	if lib, ok := r.PrivateAttr(protoc.ProtoLibraryKey).(protoc.ProtoLibrary); ok {
		return protoc.ProtoLibraryImportSpecsForKind(r.Kind(), lib)
	}
	return nil
}

// Resolve implements part of the RuleProvider interface.
func (s *ClosureJsLibrary) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
	if s.Resolver != nil {
		s.Resolver(c, ix, r, imports, from)
		return
	}

	protoc.ResolveDepsAttr("deps", false)(c, ix, r, imports, from)
	r.SetAttr("exports", r.AttrStrings("deps"))

	transitive := ResolveTransitiveProtoLibraryDeps(s.Config.Rel, r)
	descriptors := make([]string, 0)
	for _, v := range transitive {
		descriptors = append(descriptors, v)
	}

	r.SetAttr("internal_descriptors", protoc.DeduplicateAndSort(descriptors))
}

func ResolveTransitiveProtoLibraryDeps(rel string, r *rule.Rule) map[string]string {

	lib := r.PrivateAttr(protoc.ProtoLibraryKey)
	if lib == nil {
		return nil
	}
	library := lib.(protoc.ProtoLibrary)

	libRule := library.Rule()
	// already created?
	if transitiveDeps, ok := libRule.PrivateAttr(transitiveDepsKey).(map[string]string); ok {
		return transitiveDeps
	}

	// nope.
	transitiveDeps := make(map[string]string)
	resolver := protoc.GlobalResolver()

	seen := make(map[string]bool)
	stack := list.New()
	for _, src := range library.Srcs() {
		stack.PushBack(path.Join(rel, src))
	}
	// for every source file in the proto library, gather the list of source
	// files on which it depends, until there are no more unprocessed sources.
	// Foreach one check if there is an importmapping for it and record the
	// association.
	for {
		if stack.Len() == 0 {
			break
		}
		current := stack.Front()
		stack.Remove(current)

		protofile := current.Value.(string)
		if seen[protofile] {
			continue
		}
		seen[protofile] = true

		depends := resolver.Resolve("proto", "depends", protofile)
		for _, dep := range depends {
			stack.PushBack(path.Join(dep.Label.Pkg, dep.Label.Name))
		}

		result := resolver.Resolve("proto", "proto", protofile)
		if len(result) > 0 {
			first := result[0]
			transitiveDeps[protofile] = first.Label.String()
		}
	}

	libRule.SetPrivateAttr(transitiveDepsKey, transitiveDeps)

	return transitiveDeps
}
