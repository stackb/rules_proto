package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func generateMarkdown(c *Config) error {
	f, err := os.Create(c.MarkdownOut)
	if err != nil {
		return fmt.Errorf("create %s: %v", c.MarkdownOut, err)
	}
	defer f.Close()

	var buildIn, buildOut, protoFile string
	for _, src := range c.Files {
		if filepath.Ext(src) == ".proto" {
			protoFile = src
			continue
		}
		if filepath.Base(src) == "BUILD.in" {
			buildIn = src
			continue
		}
		if filepath.Base(src) == "BUILD.out" {
			buildOut = src
			continue
		}
	}

	fmt.Fprintf(f, "# %s example\n\n", filepath.Base(c.MarkdownOut))
	fmt.Fprintf(f, "\nGiven a directory with a proto file...\n\n")

	// Print the BUILD.in
	//
	if err := printFileBlock(filepath.Base(protoFile), "proto", protoFile, f); err != nil {
		return err
	}

	fmt.Fprintf(f, "\n...and a `BUILD.bazel` file with gazelle directives:\n\n")

	if err := printFileBlock("BUILD.bazel", "python", buildIn, f); err != nil {
		return err
	}

	fmt.Fprintf(f, "When gazelle is run:\n\n")

	fmt.Fprintf(f, "~~~bash\n")
	fmt.Fprintf(f, "bazel run //:gazelle\n")
	fmt.Fprintf(f, "~~~\n\n")

	fmt.Fprintf(f, "Then the following `BUILD.bazel` file will be generated:\n\n")

	// Print the BUILD.out
	//
	if err := printFileBlock("BUILD.bazel", "python", buildOut, f); err != nil {
		return err
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
	fmt.Fprintln(f, "var txtar=`")

	for _, src := range c.Files {
		dst := mapFilename(src)
		if dst == "" {
			continue
		}

		fmt.Fprintf(f, "-- %s --\n", dst)
		if dst == "WORKSPACE" {
			fmt.Fprintln(f, workspace)
			continue
		}

		data, err := ioutil.ReadFile(src)
		if err != nil {
			return fmt.Errorf("read %s: %v", src, err)
		}
		if _, err := f.Write(data); err != nil {
			return fmt.Errorf("write %s: %v", dst, err)
		}
	}

	fmt.Fprintln(f, "`")

	return nil
}

func mapFilename(in string) string {
	dir := filepath.Dir(in)
	base := filepath.Base(in)

	switch base {
	case "WORKSPACE":
		return "WORKSPACE"
	case "BUILD.in":
		return ""
	case "BUILD.out":
		return filepath.Join(dir, "BUILD.bazel")
	}

	return in
}

func printFileBlock(name, syntax, filename string, out io.Writer) error {
	fmt.Fprintf(out, "~~~%s\n", syntax)
	fmt.Fprintf(out, "# -- %s --\n", name)
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("read %s: %v", filename, err)
	}
	if _, err := out.Write(data); err != nil {
		return fmt.Errorf("write %s: %v", filename, err)
	}
	fmt.Fprintf(out, "~~~\n\n")

	return nil
}
