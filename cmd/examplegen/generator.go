package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func generateMarkdown(c *Config) error {
	f, err := os.Create(c.MarkdownOut)
	if err != nil {
		return fmt.Errorf("create %s: %v", c.MarkdownOut, err)
	}
	defer f.Close()

	var workspace, buildIn, buildOut, protoFile string
	for _, src := range c.Files {
		base := filepath.Base(src)
		ext := filepath.Ext(base)

		if ext == ".proto" {
			protoFile = src
			continue
		}

		switch base {
		case "BUILD.in":
			buildIn = src
		case "BUILD.out":
			buildOut = src
		case "MODULE.bazel":
			workspace = src
		}
	}

	if buildIn == "" {
		log.Panicf("BUILD.in not found: %+v", c)
	}
	if buildOut == "" {
		log.Panicf("BUILD.out not found: %+v", c)
	}

	fmt.Fprintf(f, "---\n")
	fmt.Fprintf(f, "layout: default\n")
	fmt.Fprintf(f, "title: %s\n", c.Name)
	fmt.Fprintf(f, "permalink: examples/%s\n", c.Name)
	fmt.Fprintf(f, "parent: Examples\n")
	fmt.Fprintf(f, "---\n\n\n")

	fmt.Fprintf(f, "# %s example\n\n", c.Name)

	fmt.Fprintf(f, "`bazel test %s_test`\n\n", c.Label)

	fmt.Fprintf(f, "\n## `BUILD.bazel` (after gazelle)\n\n")
	if err := printFileBlock("BUILD.bazel", "python", buildOut, f); err != nil {
		return err
	}

	fmt.Fprintf(f, "\n## `BUILD.bazel` (before gazelle)\n\n")
	if err := printFileBlock("BUILD.bazel", "python", buildIn, f); err != nil {
		return err
	}

	fmt.Fprintf(f, "\n## `MODULE.bazel`\n\n")
	if err := printFileBlock(filepath.Base(workspace), "python", workspace, f); err != nil {
		return err
	}

	if false {
		if err := printFileBlock(filepath.Base(protoFile), "proto", protoFile, f); err != nil {
			return err
		}
	}

	return nil
}

func generateTest(c *Config) error {
	f, err := os.Create(c.TestOut)
	if err != nil {
		return fmt.Errorf("create %s: %v", c.TestOut, err)
	}
	defer f.Close()

	fmt.Fprintln(f, testHeader)
	fmt.Fprintln(f, c.TestContent)

	fmt.Fprintln(f, "var moduleFileSuffix=`")
	if _, err := f.WriteString(c.WorkspaceIn); err != nil {
		return err
	}
	fmt.Fprintln(f, "`")

	fmt.Fprintln(f, "var txtar=`")

	for _, src := range c.Files {
		dst := mapFilename(src)
		if dst == "" {
			continue
		}

		dstFilename := filepath.Base(dst)
		if c.StripPrefix != "" {
			dstFilename = stripRel(c.StripPrefix, dst)
		}

		fmt.Fprintf(f, "-- %s --\n", dstFilename)

		data, err := os.ReadFile(src)
		if err != nil {
			return fmt.Errorf("read %q: %v", src, err)
		}
		if _, err := f.Write(data); err != nil {
			return fmt.Errorf("write %q: %v", dst, err)
		}
		f.WriteString("\n")
	}

	fmt.Fprintln(f, "`")

	return nil
}

func mapFilename(in string) string {
	dir := filepath.Dir(in)
	base := filepath.Base(in)

	switch base {
	case "MODULE.bazel":
		return ""
	case "BUILD.in":
		return ""
	case "BUILD.out":
		return filepath.Join(dir, "BUILD.bazel")
	}

	return in
}

func printFileBlock(name, syntax, filename string, out io.Writer) error {
	fmt.Fprintf(out, "~~~%s\n", syntax)
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Panicf("%s: failed to read filename=%q: %v", name, filename, err)
	}
	if _, err := out.Write(data); err != nil {
		log.Panicf("%s: write %q: %v", name, filename, err)
	}
	fmt.Fprintf(out, "~~~\n\n")

	return nil
}

// stripRel removes the rel prefix from a filename (if has matching prefix)
func stripRel(rel string, filename string) string {
	if !strings.HasPrefix(filename, rel) {
		return filename
	}
	filename = filename[len(rel):]
	return strings.TrimPrefix(filename, "/")
}
