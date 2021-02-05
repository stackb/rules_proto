package laze

import (
	"path/filepath"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/language"
)

// GenerateRules extracts build metadata from source files in a directory.
// GenerateRules is called in each directory where an update is requested
// in depth-first post-order.
//
// TODO(pcj): refactor and add proper tests for all this.
func (p *plugin) GenerateRules(args language.GenerateArgs) language.GenerateResult {
	packageName := filepath.Base(args.Dir)
	packageName = strings.ReplaceAll(packageName, "-", "_")

	pyFilenames := make([]string, 0)
	pyTestFilenames := make([]string, 0)

	// // Determine the package name by looking for the __init__.py file in the
	// // directory. If it's present, the package name is the directory name.
	// for _, f := range args.RegularFiles {
	// 	if strings.HasSuffix(f, "_test.py") {
	// 		pyTestFilenames = append(pyTestFilenames, f)
	// 	} else if strings.HasSuffix(f, ".py") {
	// 		pyFilenames = append(pyFilenames, f)
	// 	}
	// }

	// Just return nothing if the current directory has no python files.
	if len(pyFilenames) == 0 && len(pyTestFilenames) == 0 {
		return language.GenerateResult{}
	}

	// // // Collect the requirements.txt module names. If requirements.txt does
	// // // not exist, any dependency with no .py file reference will be treated
	// // // as a built-in Python module.
	// // // TODO(pcj): keep looking for the parent directories until a
	// // // requirements.txt is found. Stop at the WORKSPACE directory.
	// // requirementsTxtFilepath := filepath.Join(args.Dir, requirementsFilename)
	// // requirementsTxt, err := collectRequirementsTxt(requirementsTxtFilepath)
	// // if err != nil {
	// // 	log.Printf("ERROR: %v", err)
	// // 	return language.GenerateResult{}
	// // }

	// // libraryRules, libraryImports, err := constructPyLibraries(
	// // 	packageName,
	// // 	pyFilenames,
	// // 	args.Dir,
	// // 	requirementsTxt,
	// // )
	// // if err != nil {
	// // 	log.Printf("ERROR: %v", err)
	// // 	return language.GenerateResult{}
	// // }

	// // testRules, testImports, err := constructPyTests(
	// // 	packageName,
	// // 	pyTestFilenames,
	// // 	args.Dir,
	// // 	requirementsTxt,
	// // 	libraryRules,
	// // )
	// // if err != nil {
	// // 	log.Printf("ERROR: %v", err)
	// // 	return language.GenerateResult{}
	// // }

	// // The resulting rule targets to be generated.
	// rules := make([]*rule.Rule, 0, len(libraryRules)+len(testRules))
	// rules = append(rules, libraryRules...)
	// rules = append(rules, testRules...)

	var result language.GenerateResult
	// result.Gen = rules
	// result.Imports = append(result.Imports, libraryImports...)
	// result.Imports = append(result.Imports, testImports...)
	return result
}

// // constructPyLibraries returns the constructed py_library rule targets as well
// // as the imports.
// func constructPyLibraries(
// 	packageName string,
// 	pyFilenames []string,
// 	dir string,
// 	requirementsTxt []requirement,
// ) ([]*rule.Rule, []interface{}, error) {
// 	// Recursively collect the dependencies for all Python modules in the
// 	// current directory.
// 	packageDeps := treeset.NewWith(func(a, b interface{}) int {
// 		return godsutils.StringComparator(a.(*bzl.LiteralExpr).Token, b.(*bzl.LiteralExpr).Token)
// 	})
// 	for _, pyFilename := range pyFilenames {
// 		pyFilepath := filepath.Join(dir, pyFilename)
// 		imports, err := collectImportDependencies(pyFilepath)
// 		if err != nil {
// 			return nil, nil, fmt.Errorf("failed to construct %s: %w", pyLibraryKind, err)
// 		}

// 		// The list of non-local module names. This includes third-party AND
// 		// builtin Python modules. It is up to the target generator to determine
// 		// if the requirement should be explicitly included as a dependency or
// 		// not. This is done by comparing this list of requirements with the
// 		// user requirements in the requirements.txt file. This Gazelle
// 		// extension does not automatically update the requirements.txt file.
// 		moduleDeps := collectRequirementsFromDependencies(requirementsTxt, imports)
// 		for _, moduleDep := range moduleDeps {
// 			packageDeps.Add(moduleDep)
// 		}
// 	}

// 	pyLibrary := rule.NewRule(pyLibraryKind, packageName)
// 	pyLibrary.SetAttr("srcs", pyFilenames)

// 	packageDepsValues := packageDeps.Values()
// 	deps := make([]bzl.Expr, len(packageDepsValues))
// 	for i := range deps {
// 		deps[i] = packageDepsValues[i].(bzl.Expr)
// 	}
// 	pyLibrary.SetAttr("deps", &bzl.ListExpr{List: deps})

// 	// TODO(pcj): check for a __main__.py file to determine if a `py_binary`
// 	// should be generated as well.

// 	rules := []*rule.Rule{pyLibrary}
// 	imports := make([]interface{}, 0, len(rules))
// 	for _, rule := range rules {
// 		imports = append(imports, rule.PrivateAttr(config.GazelleImportsKey))
// 	}
// 	return rules, imports, nil
// }

// // constructPyTests returns the constructed py_test rule targets as well as the
// // imports.
// // TODO(pcj): check for a __test__.py file to determine if a `py_test`
// // should be generated as well. If not, all *_spec.py and *_test.py files should
// // be ignored.
// func constructPyTests(
// 	packageName string,
// 	pyFilenames []string,
// 	dir string,
// 	requirementsTxt []requirement,
// 	pyLibraries []*rule.Rule,
// ) ([]*rule.Rule, []interface{}, error) {
// 	// Recursively collect the dependencies for all Python modules in the
// 	// current directory.
// 	packageDeps := treeset.NewWith(func(a, b interface{}) int {
// 		return godsutils.StringComparator(a.(*bzl.LiteralExpr).Token, b.(*bzl.LiteralExpr).Token)
// 	})
// 	for _, pyFilename := range pyFilenames {
// 		pyFilepath := filepath.Join(dir, pyFilename)
// 		imports, err := collectImportDependencies(pyFilepath)
// 		if err != nil {
// 			return nil, nil, fmt.Errorf("failed to construct %s: %w", pyTestKind, err)
// 		}

// 		pyLibraryDeps := collectPyLibrariesFromDependencies(pyLibraries, imports)
// 		for _, pyLibrary := range pyLibraryDeps {
// 			packageDeps.Add(pyLibrary)
// 		}

// 		// The list of non-local module names. This includes third-party AND
// 		// builtin Python modules. It is up to the target generator to determine
// 		// if the requirement should be explicitly included as a dependency or
// 		// not. This is done by comparing this list of requirements with the
// 		// user requirements in the requirements.txt file. This Gazelle
// 		// extension does not automatically update the requirements.txt file.
// 		moduleDeps := collectRequirementsFromDependencies(requirementsTxt, imports)
// 		for _, moduleDep := range moduleDeps {
// 			packageDeps.Add(moduleDep)
// 		}
// 	}

// 	// TODO(pcj): check for testdata in the same way the Go extension for
// 	// Gazelle does.

// 	pyTest := rule.NewRule(pyTestKind, fmt.Sprintf("%s_test", packageName))
// 	pyTest.SetAttr("srcs", pyFilenames)

// 	packageDepsValues := packageDeps.Values()
// 	deps := make([]bzl.Expr, len(packageDepsValues))
// 	for i := range deps {
// 		deps[i] = packageDepsValues[i].(bzl.Expr)
// 	}
// 	pyTest.SetAttr("deps", &bzl.ListExpr{List: deps})

// 	// TODO(pcj): check for `__name__ == '__main__'` to determine if a
// 	// py_binary should be generated embedding the py_library.

// 	rules := []*rule.Rule{pyTest}
// 	imports := make([]interface{}, 0, len(rules))
// 	for _, rule := range rules {
// 		imports = append(imports, rule.PrivateAttr(config.GazelleImportsKey))
// 	}
// 	return rules, imports, nil
// }

// // collectImportDependencies collects the Python modules by recursively walking
// // the imports found on the visited modules.
// func collectImportDependencies(moduleFilepath string) ([]pyModule, error) {
// 	result := make([]pyModule, 0)

// 	parser := &python3ModuleParser{}
// 	if err := parser.parse(moduleFilepath); err != nil {
// 		return nil, fmt.Errorf("failed to collect import dependencies from %q: %w", moduleFilepath, err)
// 	}

// 	moduleDir := filepath.Dir(moduleFilepath)

// 	for _, name := range parser.dependencyNames() {
// 		m := pyModule{name: name}

// 		// If the import is in a local module file, process it recursively.
// 		submoduleFilepath := filepath.Join(moduleDir, fmt.Sprintf("%s.py", name))
// 		if _, err := os.Stat(submoduleFilepath); err == nil {
// 			m.filepath = submoduleFilepath
// 			// Recurse
// 			submoduleDependencies, err := collectImportDependencies(submoduleFilepath)
// 			if err != nil {
// 				return nil, fmt.Errorf("failed to collect import dependencies from %q: %w", moduleFilepath, err)
// 			}
// 			m.dependencies = submoduleDependencies
// 		}

// 		result = append(result, m)
// 	}

// 	return result, nil
// }

// // collectRequirementsTxt collects the list of requirements from a
// // requirements.txt file.
// func collectRequirementsTxt(requirementsTxtFilepath string) ([]requirement, error) {
// 	result := make([]requirement, 0)

// 	if _, err := os.Stat(requirementsTxtFilepath); err == nil {
// 		file, err := os.Open(requirementsTxtFilepath)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to collect requirements from %q: %w", requirementsTxtFilepath, err)
// 		}

// 		var currentRequirement requirement
// 		var continuePreviousLine bool

// 		scanner := bufio.NewScanner(file)
// 		for scanner.Scan() {
// 			line := scanner.Text()

// 			// Ignore empty lines.
// 			if len(line) == 0 {
// 				continue
// 			}

// 			// Ignore lines containing only comments.
// 			if strings.HasPrefix(line, "#") {
// 				continue
// 			}

// 			// If this line is a continuation of the previous one, append the
// 			// current line to the current requirement being processed,
// 			// otherwise, start a new requirement.
// 			if continuePreviousLine {
// 				currentRequirement.raw += strings.TrimSuffix(strings.TrimSpace(line), `\`)
// 			} else {
// 				currentRequirement = requirement{}
// 				currentRequirement.raw = strings.TrimSuffix(strings.TrimSpace(line), `\`)
// 			}

// 			// Control whether the next line in the requirements.txt should be
// 			// a continuation of the current requirement being processed or not.
// 			continuePreviousLine = strings.HasSuffix(line, `\`)

// 			if !continuePreviousLine {
// 				currentRequirement.parse()
// 				result = append(result, currentRequirement)
// 			}
// 		}

// 		if err := scanner.Err(); err != nil {
// 			return nil, fmt.Errorf("failed to collect requirements from %q: %w", requirementsTxtFilepath, err)
// 		}
// 	}

// 	return result, nil
// }

// // collectPyLibrariesFromDependencies collects the py_library targets that are
// // found matching in the given module dependencies and constructs a list of
// // bzl.Expr for them.
// func collectPyLibrariesFromDependencies(pyLibraries []*rule.Rule, dependencies []pyModule) []bzl.Expr {
// 	result := make([]bzl.Expr, 0)
// 	for _, dependency := range dependencies {
// 		for _, pyLibrary := range pyLibraries {
// 			for _, filename := range pyLibrary.AttrStrings("srcs") {
// 				moduleName := strings.TrimSuffix(filename, ".py")
// 				if dependency.name == moduleName {
// 					result = append(result, &bzl.LiteralExpr{Token: fmt.Sprintf(`":%s"`, pyLibrary.Name())})
// 					break
// 				}
// 			}
// 		}
// 		if len(dependency.dependencies) > 0 {
// 			result = append(result, collectPyLibrariesFromDependencies(pyLibraries, dependency.dependencies)...)
// 		}
// 	}
// 	return result
// }

// // collectRequirementsFromDependencies collects the requirements from the parsed
// // requirements.txt that are found matching in the given module dependencies and
// // constructs a list of bzl.Expr for them.
// // TODO(pcj): not having a list of builtins can be problematic on a dynamic
// // language like Python because an error to generate a requirement("") can
// // result in a runtime error. How to make this more robust?
// func collectRequirementsFromDependencies(reqs []requirement, dependencies []pyModule) []bzl.Expr {
// 	result := make([]bzl.Expr, 0)
// 	for _, dependency := range dependencies {
// 		for _, requirement := range reqs {
// 			if requirement.name == dependency.name {
// 				result = append(result, &bzl.LiteralExpr{Token: fmt.Sprintf(`requirement("%s")`, requirement.name)})
// 			}
// 		}
// 		if len(dependency.dependencies) > 0 {
// 			result = append(result, collectRequirementsFromDependencies(reqs, dependency.dependencies)...)
// 		}
// 	}
// 	return result
// }
