package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

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
	BuildFile        string
	Name             string
	Path             string
	Label            string
	RepositoryRule   string
	Sha256           string
	StripPrefix      string
	Urls             []string
	WorkspaceSnippet string
	Deps    []*ProtoDependencyInfo
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

type dependency struct {
	ParentName string
	Dep        *ProtoDependencyInfo
}
