/* Copyright 2017 The Bazel Authors. All rights reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/bazelbuild/bazel-gazelle/config"
	gzflag "github.com/bazelbuild/bazel-gazelle/flag"
	"github.com/bazelbuild/bazel-gazelle/label"
	"github.com/bazelbuild/bazel-gazelle/language"
	"github.com/bazelbuild/bazel-gazelle/merger"
	"github.com/bazelbuild/bazel-gazelle/repo"
	"github.com/bazelbuild/bazel-gazelle/resolve"
	"github.com/bazelbuild/bazel-gazelle/rule"
	"github.com/bazelbuild/bazel-gazelle/walk"
)

// updateConfig holds configuration information needed to run the fix and
// update commands. This includes everything in config.Config, but it also
// includes some additional fields that aren't relevant to other packages.
type updateConfig struct {
	dirs     []string
	emit     emitFunc
	repos    []repo.Repo
	walkMode walk.Mode
}

type emitFunc func(c *config.Config, f *rule.File) error

const updateName = "_update"

func getUpdateConfig(c *config.Config) *updateConfig {
	return c.Exts[updateName].(*updateConfig)
}

type updateConfigurer struct {
	recursive    bool
	knownImports []string
}

func (ucr *updateConfigurer) RegisterFlags(fs *flag.FlagSet, cmd string, c *config.Config) {
	uc := &updateConfig{}
	c.Exts[updateName] = uc

	c.ShouldFix = cmd == "fix"

	fs.BoolVar(&ucr.recursive, "r", true, "when true, gazelle will update subdirectories recursively")
	fs.Var(&gzflag.MultiFlag{Values: &ucr.knownImports}, "known_import", "import path for which external resolution is skipped (can specify multiple times)")
}

func (ucr *updateConfigurer) CheckFlags(fs *flag.FlagSet, c *config.Config) error {
	uc := getUpdateConfig(c)
	uc.emit = fixFile

	dirs := fs.Args()
	if len(dirs) == 0 {
		dirs = []string{"."}
	}
	uc.dirs = make([]string, len(dirs))
	for i, arg := range dirs {
		dir := arg
		if !filepath.IsAbs(dir) {
			dir = filepath.Join(c.WorkDir, dir)
		}
		dir, err := filepath.EvalSymlinks(dir)
		if err != nil {
			return fmt.Errorf("%s: failed to resolve symlinks: %v", arg, err)
		}
		if !isDescendingDir(dir, c.RepoRoot) {
			return fmt.Errorf("%s: not a subdirectory of repo root %s", arg, c.RepoRoot)
		}
		uc.dirs[i] = dir
	}

	if ucr.recursive && c.IndexLibraries {
		uc.walkMode = walk.VisitAllUpdateSubdirsMode
	} else if c.IndexLibraries {
		uc.walkMode = walk.VisitAllUpdateDirsMode
	} else if ucr.recursive {
		uc.walkMode = walk.UpdateSubdirsMode
	} else {
		uc.walkMode = walk.UpdateDirsMode
	}

	return nil
}

func (ucr *updateConfigurer) KnownDirectives() []string { return nil }

func (ucr *updateConfigurer) Configure(c *config.Config, rel string, f *rule.File) {}

// visitRecord stores information about about a directory visited with
// packages.Walk.
type visitRecord struct {
	// pkgRel is the slash-separated path to the visited directory, relative to
	// the repository root. "" for the repository root itself.
	pkgRel string

	// c is the configuration for the directory with directives applied.
	c *config.Config

	// rules is a list of generated Go rules.
	rules []*rule.Rule

	// imports contains opaque import information for each rule in rules.
	imports []interface{}

	// empty is a list of empty Go rules that may be deleted.
	empty []*rule.Rule

	// file is the build file being processed.
	file *rule.File

	// mappedKinds are mapped kinds used during this visit.
	mappedKinds    []config.MappedKind
	mappedKindInfo map[string]rule.KindInfo
}

var genericLoads = []rule.LoadInfo{
	{
		Name:    "@bazel_gazelle//:def.bzl",
		Symbols: []string{"gazelle"},
	},
}

func runFixUpdate(wd string, cmd command, args []string) (err error) {
	cexts := make([]config.Configurer, 0, len(languages)+3)
	cexts = append(cexts,
		&config.CommonConfigurer{},
		&updateConfigurer{},
		&walk.Configurer{},
		&resolve.Configurer{})
	mrslv := newMetaResolver()
	kinds := make(map[string]rule.KindInfo)
	loads := genericLoads
	exts := make([]interface{}, 0, len(languages))
	for _, lang := range languages {
		cexts = append(cexts, lang)
		for kind, info := range lang.Kinds() {
			mrslv.AddBuiltin(kind, lang)
			kinds[kind] = info
		}
		loads = append(loads, lang.Loads()...)
		exts = append(exts, lang)
	}
	ruleIndex := resolve.NewRuleIndex(mrslv.Resolver, exts...)

	c, err := newFixUpdateConfiguration(wd, cmd, args, cexts)
	if err != nil {
		return err
	}

	// Visit all directories in the repository.
	var visits []visitRecord
	uc := getUpdateConfig(c)

	walk.Walk(c, cexts, uc.dirs, uc.walkMode, func(dir, rel string, c *config.Config, update bool, f *rule.File, subdirs, regularFiles, genFiles []string) {
		// If this file is ignored or if Gazelle was not asked to update this
		// directory, just index the build file and move on.
		if !update {
			if c.IndexLibraries && f != nil {
				for _, r := range f.Rules {
					ruleIndex.AddRule(c, r, f)
				}
			}
			return
		}

		// Fix any problems in the file.
		if f != nil {
			for _, l := range filterLanguages(c, languages) {
				l.Fix(c, f)
			}
		}

		// Generate rules.
		var empty, gen []*rule.Rule
		var imports []interface{}
		for _, l := range filterLanguages(c, languages) {
			res := l.GenerateRules(language.GenerateArgs{
				Config:       c,
				Dir:          dir,
				Rel:          rel,
				File:         f,
				Subdirs:      subdirs,
				RegularFiles: regularFiles,
				GenFiles:     genFiles,
				OtherEmpty:   empty,
				OtherGen:     gen})
			if len(res.Gen) != len(res.Imports) {
				log.Panicf("%s: language %s generated %d rules but returned %d imports", rel, l.Name(), len(res.Gen), len(res.Imports))
			}
			empty = append(empty, res.Empty...)
			gen = append(gen, res.Gen...)
			imports = append(imports, res.Imports...)
		}
		if f == nil && len(gen) == 0 {
			return
		}

		// Apply and record relevant kind mappings.
		var (
			mappedKinds    []config.MappedKind
			mappedKindInfo = make(map[string]rule.KindInfo)
		)
		for _, r := range gen {
			if repl, ok := c.KindMap[r.Kind()]; ok {
				mappedKindInfo[repl.KindName] = kinds[r.Kind()]
				mappedKinds = append(mappedKinds, repl)
				mrslv.MappedKind(rel, repl)
				r.SetKind(repl.KindName)
			}
		}

		// Insert or merge rules into the build file.
		if f == nil {
			f = rule.EmptyFile(filepath.Join(dir, c.DefaultBuildFileName()), rel)
			for _, r := range gen {
				r.Insert(f)
			}
		} else {
			merger.MergeFile(f, empty, gen, merger.PreResolve,
				unionKindInfoMaps(kinds, mappedKindInfo))
		}
		visits = append(visits, visitRecord{
			pkgRel:         rel,
			c:              c,
			rules:          gen,
			imports:        imports,
			empty:          empty,
			file:           f,
			mappedKinds:    mappedKinds,
			mappedKindInfo: mappedKindInfo,
		})

		// Add library rules to the dependency resolution table.
		if c.IndexLibraries {
			for _, r := range f.Rules {
				ruleIndex.AddRule(c, r, f)
			}
		}
	})

	// Finish building the index for dependency resolution.
	ruleIndex.Finish()

	// Resolve dependencies.
	rc, cleanupRc := repo.NewRemoteCache(uc.repos)
	defer func() {
		if cerr := cleanupRc(); err == nil && cerr != nil {
			err = cerr
		}
	}()
	for _, v := range visits {
		for i, r := range v.rules {
			from := label.New(c.RepoName, v.pkgRel, r.Name())
			if rslv := mrslv.Resolver(r, v.pkgRel); rslv != nil {
				rslv.Resolve(v.c, ruleIndex, rc, r, v.imports[i], from)
			}
		}
		merger.MergeFile(v.file, v.empty, v.rules, merger.PostResolve,
			unionKindInfoMaps(kinds, v.mappedKindInfo))
	}

	// Emit merged files.
	var exit error
	for _, v := range visits {
		merger.FixLoads(v.file, applyKindMappings(v.mappedKinds, loads))
		if err := uc.emit(v.c, v.file); err != nil {
			if err == exitError {
				exit = err
			} else {
				log.Print(err)
			}
		}
	}

	return exit
}

func newFixUpdateConfiguration(wd string, cmd command, args []string, cexts []config.Configurer) (*config.Config, error) {
	c := config.New()
	c.WorkDir = wd

	fs := flag.NewFlagSet("gazelle", flag.ContinueOnError)
	// Flag will call this on any parse error. Don't print usage unless
	// -h or -help were passed explicitly.
	fs.Usage = func() {}

	for _, cext := range cexts {
		cext.RegisterFlags(fs, cmd.String(), c)
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			fixUpdateUsage(fs)
			return nil, err
		}
		// flag already prints the error; don't print it again.
		log.Fatal("Try -help for more information.")
	}

	for _, cext := range cexts {
		if err := cext.CheckFlags(fs, c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func fixUpdateUsage(fs *flag.FlagSet) {
	fmt.Fprint(os.Stderr, `usage: gazelle [fix|update] [flags...] [package-dirs...]

The update command creates new build files and update existing BUILD files
when needed.

The fix command also creates and updates build files, and in addition, it may
make potentially breaking updates to usage of rules. For example, it may
delete obsolete rules or rename existing rules.

There are several output modes which can be selected with the -mode flag. The
output mode determines what Gazelle does with updated BUILD files.

  fix (default) - write updated BUILD files back to disk.

Gazelle accepts a list of paths to Go package directories to process (defaults
to the working directory if none are given). It recursively traverses
subdirectories. All directories must be under the directory specified by
-repo_root; if -repo_root is not given, this is the directory containing the
WORKSPACE file.

FLAGS:

`)
	fs.PrintDefaults()
}

func isDescendingDir(dir, root string) bool {
	rel, err := filepath.Rel(root, dir)
	if err != nil {
		return false
	}
	if rel == "." {
		return true
	}
	return !strings.HasPrefix(rel, "..")
}

func findOutputPath(c *config.Config, f *rule.File) string {
	if c.ReadBuildFilesDir == "" && c.WriteBuildFilesDir == "" {
		return f.Path
	}
	baseDir := c.WriteBuildFilesDir
	if c.WriteBuildFilesDir == "" {
		baseDir = c.RepoRoot
	}
	outputDir := filepath.Join(baseDir, filepath.FromSlash(f.Pkg))
	defaultOutputPath := filepath.Join(outputDir, c.DefaultBuildFileName())
	files, err := ioutil.ReadDir(outputDir)
	if err != nil {
		// Ignore error. Directory probably doesn't exist.
		return defaultOutputPath
	}
	outputPath := rule.MatchBuildFileName(outputDir, c.ValidBuildFileNames, files)
	if outputPath == "" {
		return defaultOutputPath
	}
	return outputPath
}

func unionKindInfoMaps(a, b map[string]rule.KindInfo) map[string]rule.KindInfo {
	if len(a) == 0 {
		return b
	}
	if len(b) == 0 {
		return a
	}
	result := make(map[string]rule.KindInfo, len(a)+len(b))
	for _, m := range []map[string]rule.KindInfo{a, b} {
		for k, v := range m {
			result[k] = v
		}
	}
	return result
}

// applyKindMappings returns a copy of LoadInfo that includes c.KindMap.
func applyKindMappings(mappedKinds []config.MappedKind, loads []rule.LoadInfo) []rule.LoadInfo {
	if len(mappedKinds) == 0 {
		return loads
	}

	// Add new RuleInfos or replace existing ones with merged ones.
	mappedLoads := make([]rule.LoadInfo, len(loads))
	copy(mappedLoads, loads)
	for _, mappedKind := range mappedKinds {
		mappedLoads = appendOrMergeKindMapping(mappedLoads, mappedKind)
	}
	return mappedLoads
}

// appendOrMergeKindMapping adds LoadInfo for the given replacement.
func appendOrMergeKindMapping(mappedLoads []rule.LoadInfo, mappedKind config.MappedKind) []rule.LoadInfo {
	// If mappedKind.KindLoad already exists in the list, create a merged copy.
	for i, load := range mappedLoads {
		if load.Name == mappedKind.KindLoad {
			mappedLoads[i].Symbols = append(load.Symbols, mappedKind.KindName)
			return mappedLoads
		}
	}

	// Add a new LoadInfo.
	return append(mappedLoads, rule.LoadInfo{
		Name:    mappedKind.KindLoad,
		Symbols: []string{mappedKind.KindName},
	})
}
