package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func generate(c *Config) error {
	f, err := os.Create(c.Out)
	if err != nil {
		return fmt.Errorf("create %s: %v", c.Out, err)
	}
	defer f.Close()

	fmt.Fprintln(f, header)
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
