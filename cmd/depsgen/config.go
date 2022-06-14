package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// knownLoads is a mapping between load symbols and their source.  Those with an
// empty string are native and do not need a load statement.
var knownLoads = map[string]string{
	"bind":             "",
	"go_repository":    "@bazel_gazelle//:deps.bzl",
	"http_archive":     "@bazel_tools//tools/build_defs/repo:http.bzl",
	"http_file":        "@bazel_tools//tools/build_defs/repo:http.bzl",
	"local_repository": "",
	"npm_install":      "@build_bazel_rules_nodejs//:index.bzl",
	"yarn_install":     "@build_bazel_rules_nodejs//:index.bzl",
}

type Config struct {
	// Out is the filename to write
	Out string
	// Name is a prefix for the generated deps macro function.
	Name string
	// Deps is the list of ProtoDependencyInfos
	Deps []*ProtoDependencyInfo
}

// ProtoDependencyInfo represents the starlark ProtoDependencyInfo provider.
// The fields are a mashup of all possible fields from repository rules.
type ProtoDependencyInfo struct {
	BuildFile          string
	BuildFileContent   string
	BuildFileProtoMode string
	Deps               []*ProtoDependencyInfo
	Executable         bool
	Importpath         string
	Label              string
	Name               string
	Path               string
	RepositoryRule     string
	Sha256             string
	StripPrefix        string
	SymlinkNodeModules bool
	FrozenLockfile     bool
	Sum                string
	Urls               []string
	Version            string
	WorkspaceSnippet   string
	PackageJson        string
	PackageLockJson    string
	PatchArgs          []string
	Patches            []string
	YarnLock           string
}

type LoadInfo struct {
	Label   string
	Symbols []string
}

// fromJSON constructs a Config struct from the given filename that contains a
// JSON.
func fromJSON(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("unmarshal: %w", err)
	}

	return &config, nil
}

// collectDeps accumulates the transitive dependencies.
func collectDeps(top []*ProtoDependencyInfo) (deps []*dependency) {
	seen := make(map[string]bool)

	var visit func(string, *ProtoDependencyInfo)
	visit = func(parentName string, dep *ProtoDependencyInfo) {
		if seen[dep.Name] {
			return
		}
		seen[dep.Name] = true
		for _, child := range dep.Deps {
			visit(dep.Name, child)
		}
		deps = append(deps, &dependency{
			ParentName: parentName,
			Dep:        dep,
		})
	}

	for _, dep := range top {
		if seen[dep.Name] {
			continue
		}
		visit("<TOP>", dep)
	}

	return
}

// collectLoads accumulates the required load statements.
func collectLoads(deps []*dependency) []*LoadInfo {
	required := make(map[string][]string)

	for _, item := range deps {
		symbol := item.Dep.RepositoryRule
		source := knownLoads[symbol]
		if source == "" {
			continue
		}
		required[source] = append(required[source], symbol)
	}

	loads := make([]*LoadInfo, 0)

	for source, symbols := range required {
		load := &LoadInfo{Label: source}
		loads = append(loads, load)
		seen := make(map[string]bool)
		for _, symbol := range symbols {
			if seen[symbol] {
				continue
			}
			seen[symbol] = true
			load.Symbols = append(load.Symbols, symbol)
		}
	}

	return loads
}

type dependency struct {
	ParentName string
	Dep        *ProtoDependencyInfo
}
