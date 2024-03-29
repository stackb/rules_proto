package protobuf

import (
	"log"
	"sort"

	"github.com/bazelbuild/bazel-gazelle/rule"

	"github.com/stackb/rules_proto/pkg/protoc"
)

// Kinds returns a map of maps rule names (kinds) and information on how to
// match and merge attributes that may be found in rules of those kinds. All
// kinds of rules generated for this language may be found here.
func (*protobufLang) Kinds() map[string]rule.KindInfo {
	registry := protoc.Rules()

	kinds := make(map[string]rule.KindInfo)
	kinds[overrideKindName] = overrideKind

	for _, name := range registry.RuleNames() {
		rule, err := registry.LookupRule(name)
		if err != nil {
			log.Fatal("Kinds:", err)
		}
		if _, ok := kinds[rule.Name()]; ok {
			log.Fatal("Kinds: duplicate rule name:", rule.Name())
		}
		kinds[rule.Name()] = rule.KindInfo()
	}

	return kinds
}

// Loads returns .bzl files and symbols they define. Every rule generated by
// GenerateRules, now or in the past, should be loadable from one of these
// files.
func (pl *protobufLang) Loads() []rule.LoadInfo {

	// Merge symbols
	symbolsByLoadName := make(map[string][]string)
	for _, name := range pl.rules.RuleNames() {
		rule, err := pl.rules.LookupRule(name)
		if err != nil {
			log.Fatal(err)
		}
		load := rule.LoadInfo()
		if load.Name == "" {
			log.Fatal("Loads: empty load name for rule:", name)
		}
		symbolsByLoadName[load.Name] = append(symbolsByLoadName[load.Name], load.Symbols...)
	}

	// Ensure names are sorted otherwise order of load statements can be
	// non-deterministic
	keys := make([]string, 0)
	for name := range symbolsByLoadName {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	// Build final load list
	loads := make([]rule.LoadInfo, 0)
	for _, name := range keys {
		symbols := symbolsByLoadName[name]
		sort.Strings(symbols)
		loads = append(loads, rule.LoadInfo{
			Name:    name,
			Symbols: symbols,
		})
	}

	return loads
}
