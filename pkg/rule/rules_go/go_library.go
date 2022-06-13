package rules_go

import (
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

const (
	ProtoGoLibraryRuleName = "proto_go_library"
	goLibraryRuleSuffix    = "_go_proto"
)

func init() {
	protoc.Rules().MustRegisterRule("stackb:rules_proto:"+ProtoGoLibraryRuleName,
		&goLibrary{
			kindName:             ProtoGoLibraryRuleName,
			protoLibrariesByRule: make(map[label.Label][]protoc.ProtoLibrary),
		})
}

// goLibrary implements LanguageRule for the '{proto|grpc}_go_library' rule from
// @rules_proto.
type goLibrary struct {
	kindName             string
	protoLibrariesByRule map[label.Label][]protoc.ProtoLibrary
}

// Name implements part of the LanguageRule interface.
func (s *goLibrary) Name() string {
	return s.kindName
}

// KindInfo implements part of the LanguageRule interface.
func (s *goLibrary) KindInfo() rule.KindInfo {
	return rule.KindInfo{
		MatchAttrs: []string{"importpath"},
		NonEmptyAttrs: map[string]bool{
			"deps":  true,
			"embed": true,
			"srcs":  true,
		},
		MergeableAttrs: map[string]bool{
			"embed": true,
			"srcs":  true,
		},
		ResolveAttrs: map[string]bool{"deps": true},
	}
}

// LoadInfo implements part of the LanguageRule interface.
func (s *goLibrary) LoadInfo() rule.LoadInfo {
	return rule.LoadInfo{
		Name:    fmt.Sprintf("@build_stack_rules_proto//rules/go:%s.bzl", s.kindName),
		Symbols: []string{s.kindName},
	}
}

// ProvideRule implements part of the LanguageRule interface.
func (s *goLibrary) ProvideRule(cfg *protoc.LanguageRuleConfig, pc *protoc.ProtocConfiguration) protoc.RuleProvider {
	// collect the outputs and the deps.  Search all the PluginConfigurations.
	// If the produced .go files, include them and add their deps.
	outputs := make([]string, 0)
	pluginDeps := make([]string, 0)

	for _, pluginConfig := range pc.Plugins {
		for _, out := range pluginConfig.Outputs {
			if path.Ext(out) == ".go" {
				outputs = append(outputs, out)
				pluginDeps = append(pluginDeps, pluginConfig.Config.GetDeps()...)
			}
		}
	}

	if len(outputs) == 0 {
		return nil
	}

	for i, output := range outputs {
		outputs[i] = path.Join(pc.Rel, path.Base(output))
	}

	rule := &goLibraryRule{
		kindName:             s.kindName,
		ruleNameSuffix:       goLibraryRuleSuffix,
		outputs:              protoc.DeduplicateAndSort(outputs),
		deps:                 protoc.DeduplicateAndSort(pluginDeps),
		ruleConfig:           cfg,
		pc:                   pc,
		protoLibrariesByRule: s.protoLibrariesByRule,
	}
	rule.id = label.New("", pc.Rel, rule.Name())
	return rule
}

// goLibraryRule implements RuleProvider for 'go_library'-derived rules.
type goLibraryRule struct {
	id                   label.Label
	kindName             string
	ruleNameSuffix       string
	outputs              []string
	deps                 []string
	pc                   *protoc.ProtocConfiguration
	ruleConfig           *protoc.LanguageRuleConfig
	protoLibrariesByRule map[label.Label][]protoc.ProtoLibrary
}

// Kind implements part of the ruleProvider interface.
func (s *goLibraryRule) Kind() string {
	return s.kindName
}

// Name implements part of the ruleProvider interface.
func (s *goLibraryRule) Name() string {
	return s.pc.Library.BaseName() + s.ruleNameSuffix
}

// Srcs computes the srcs list for the rule.
func (s *goLibraryRule) Srcs() []string {
	srcs := make([]string, 0)
	for _, output := range s.outputs {
		srcs = append(srcs, protoc.StripRel(s.pc.Rel, output))
	}
	return srcs
}

// deps returns all known "configured" dependencies:
// 1. Those given by the plugin implementations that contributed outputs (their 'deps' directive).
// 2. Those given by 'deps' directive on the rule config.
// 3. Those given by resolving "rewrite" specs against the proto file imports.
func (s *goLibraryRule) configDeps() []string {
	deps := s.deps
	deps = append(deps, s.ruleConfig.GetDeps()...)
	resolvedDeps := protoc.ResolveLibraryRewrites(s.ruleConfig.GetRewrites(), s.pc.Library)
	deps = append(deps, resolvedDeps...)
	return deps
}

// Visibility provides visibility labels.
func (s *goLibraryRule) Visibility() []string {
	return s.ruleConfig.GetVisibility()
}

func (s *goLibraryRule) goPrefix() string {
	res := protoc.GlobalResolver().Resolve("gazelle", "directive", "prefix")
	if len(res) == 0 {
		return ""
	}
	return res[0].Label.Pkg
}

// importPath computes the import path.
func (s *goLibraryRule) importPath() string {
	// Try 'M' options first
	if imp := s.getPluginImportMappingOption(); imp != "" {
		return imp
	}

	// Next try the 'go_package' option in an imported file
	for _, file := range s.pc.Library.Files() {
		if value, _ := protoc.GetNamedOption(file.Options(), "go_package"); value != "" {
			if strings.LastIndexByte(value, '/') == -1 {
				// return langgo.InferImportPath(c, rel)
				continue // TODO: do more research here on if this is the correct approach
			}
			if i := strings.LastIndexByte(value, ';'); i != -1 {
				return value[:i]
			}
			return value
		}
	}

	// fallback to 'gazelle:prefix + rel'
	prefix := s.goPrefix()
	if prefix == "" {
		return ""
	}

	pkg := s.pc.Rel
	name := s.pc.Library.BaseName()

	return path.Join(prefix, pkg, name)
}

// Rule implements part of the ruleProvider interface.
func (s *goLibraryRule) Rule(otherGen ...*rule.Rule) *rule.Rule {
	importpath := s.importPath()
	srcs := s.Srcs()
	deps := s.configDeps()
	visibility := s.Visibility()
	imports := s.pc.Library.Imports()

	// Check if an existing proto_go_library rule has already been generated
	// under this importpath.  If so, we need to merge into it rather than
	// create a new rule.
	for _, other := range otherGen {
		if other.Kind() == ProtoGoLibraryRuleName && other.AttrString("importpath") == importpath {
			otherLabel := label.New("", s.pc.Rel, other.Name())
			otherSrcs := other.AttrStrings("srcs")
			otherDeps := other.AttrStrings("deps")
			otherVis := other.AttrStrings("visibility")
			otherImports := other.PrivateAttr(config.GazelleImportsKey).([]string)

			other.SetAttr("srcs", protoc.DeduplicateAndSort(append(otherSrcs, srcs...)))
			other.SetAttr("deps", protoc.DeduplicateAndSort(append(otherDeps, deps...)))
			other.SetAttr("visibility", protoc.DeduplicateAndSort(append(otherVis, visibility...)))
			other.SetPrivateAttr(config.GazelleImportsKey, protoc.DeduplicateAndSort(append(otherImports, imports...)))

			s.protoLibrariesByRule[otherLabel] = append(s.protoLibrariesByRule[otherLabel], s.pc.Library)

			return other
		}
	}

	newRule := rule.NewRule(s.Kind(), s.Name())
	newRule.SetAttr("srcs", srcs)
	newRule.SetPrivateAttr(config.GazelleImportsKey, imports)
	s.protoLibrariesByRule[s.id] = []protoc.ProtoLibrary{s.pc.Library}

	if importpath != "" {
		newRule.SetAttr("importpath", importpath)
	}
	if len(deps) > 0 {
		newRule.SetAttr("deps", deps)
	}
	if len(visibility) > 0 {
		newRule.SetAttr("visibility", visibility)
	}
	return newRule
}

func printProtoLibraryNames(libs []protoc.ProtoLibrary) string {
	names := make([]string, len(libs))
	for i, lib := range libs {
		names[i] = lib.BaseName()
	}
	return strings.Join(names, ",")
}

// Imports implements part of the RuleProvider interface.
func (s *goLibraryRule) Imports(c *config.Config, r *rule.Rule, f *rule.File) []resolve.ImportSpec {
	// for the cross-resolver such that go can cross-resolve this library
	from := label.New("", f.Pkg, r.Name())

	// log.Println("provide for cross-resolver", r.AttrString("importpath"), from)
	protoc.GlobalResolver().Provide("go", "go", r.AttrString("importpath"), from)

	libs, ok := s.protoLibrariesByRule[s.id]
	if !ok {
		return nil
	}

	return protoc.ProtoLibraryImportSpecsForKind(r.Kind(), libs...)
}

// Resolve implements part of the RuleProvider interface.
func (s *goLibraryRule) Resolve(c *config.Config, ix *resolve.RuleIndex, r *rule.Rule, imports []string, from label.Label) {
	protoc.ResolveDepsAttr("deps", true)(c, ix, r, imports, from)

	// need to make one more pass to possibly move deps into embeds.  There may
	// be dependencies *IN OTHER PACKAGES* that have the same importpath; in
	// that case we need to embed, not depend.
	all := r.AttrStrings("deps")

	deps := make([]string, 0)
	embeds := make([]string, 0)
	importpath := r.AttrString("importpath")

	for _, dep := range all {
		dl, err := label.Parse(dep)
		if err != nil {
			log.Printf("resolve deps failed for for rule %s %s: label parse %q error: %v", r.Kind(), r.Name(), dep, err)
			deps = append(deps, dep)
			continue
		}

		// If this is a relative label, make it absolute
		if dl.Repo == "" && dl.Pkg == "" {
			dl = label.Label{Pkg: s.pc.Rel, Name: dl.Name}
		}

		// retrieve the rule for this label
		if dr := protoc.GlobalRuleIndex().Get(dl); dr != nil {
			depImportpath := dr.AttrString("importpath")
			// if it has the same importpath, need to embed this
			if depImportpath == importpath {
				embeds = append(embeds, dep)
				continue
			}
		}

		deps = append(deps, dep)
	}

	if len(deps) > 0 {
		r.SetAttr("deps", protoc.DeduplicateAndSort(deps))
	}
	if len(embeds) > 0 {
		r.SetAttr("embed", protoc.DeduplicateAndSort(embeds))
	}
}

func (s *goLibraryRule) getPluginImportMappingOption() string {
	// first, iterate all the plugins and gather options that look like
	// protoc-gen-go "importmapping" (M) options (e.g
	// Mfoo.proto=github.com/example/foo).
	mappings := make(map[string]string)

	tryParseMapping := func(opt string) {
		if !strings.HasPrefix(opt, "M") {
			return
		}
		parts := strings.SplitN(opt[1:], "=", 2)
		if len(parts) != 2 {
			return
		}
		mappings[parts[0]] = parts[1]
	}

	// search all plugins
	for _, plugin := range s.pc.Plugins {
		for _, opt := range plugin.Options {
			tryParseMapping(opt)
		}
	}
	// and all rule options
	for _, opt := range s.ruleConfig.GetOptions() {
		tryParseMapping(opt)
	}

	// now that we've gathered all possible options; search all library files
	// (e.g. foo.proto) and see if we can find a match.
	for _, file := range s.pc.Library.Files() {
		filename := path.Join(file.Dir, file.Basename)
		mapping := mappings[filename]
		if mapping != "" {
			return mapping
		}
	}

	return ""
}
