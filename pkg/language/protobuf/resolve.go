package protobuf

import (
	"log"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

// Imports returns a list of ImportSpecs that can be used to import the rule r.
// This is used to populate RuleIndex.
//
// If nil is returned, the rule will not be indexed. If any non-nil slice is
// returned, including an empty slice, the rule will be indexed.
func (pl *protobufLang) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	if resolver, ok := r.PrivateAttr(protoc.RuleProviderKey).(protoc.RuleProvider); ok {
		imports := resolver.Imports(c, r, f)
		from := label.New("", f.Pkg, r.Name())
		for _, imp := range imports {
			protoc.GlobalResolver().Provide(
				pl.name,
				imp.Lang,
				imp.Imp,
				from,
			)
		}
	}
	return nil
}

// Embeds returns a list of labels of rules that the given rule embeds. If a
// rule is embedded by another importable rule of the same language, only the
// embedding rule will be indexed. The embedding rule will inherit the imports
// of the embedded rule. Since SkyLark doesn't support embedding this should
// always return nil.
func (*protobufLang) Embeds(r *rule.Rule, from label.Label) []label.Label { return nil }

// Resolve translates imported libraries for a given rule into Bazel
// dependencies. Information about imported libraries is returned for each rule
// generated by language.GenerateRules in language.GenerateResult.Imports.
// Resolve generates a "deps" attribute (or the appropriate language-specific
// equivalent) for each import according to language-specific rules and
// heuristics.
func (pl *protobufLang) Resolve(
	c *config.Config,
	ix *resolve.RuleIndex,
	rc *repo.RemoteCache,
	r *rule.Rule,
	importsRaw interface{},
	from label.Label,
) {
	if r.Kind() == overrideKindName {
		resolveOverrideRule(from.Pkg, r)
		return
	}

	if resolver, ok := r.PrivateAttr(protoc.RuleProviderKey).(protoc.RuleProvider); ok {
		if imports, ok := importsRaw.([]string); ok {
			resolver.Resolve(c, ix, r, imports, from)
		} else {
			log.Panicf("warning: resolve imports: expected []string, got %T", importsRaw)
		}
	}
}

func (*protobufLang) CrossResolve(c *config.Config, ix *resolve.RuleIndex, imp resolve.ImportSpec, lang string) []resolve.FindResult {
	return protoc.GlobalResolver().CrossResolve(c, ix, imp, lang)
}
