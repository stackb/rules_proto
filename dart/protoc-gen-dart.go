package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"syscall"
)

func main() {
	tempDir, err := ioutil.TempDir("", "protoc-gen-dart-")
	defer os.RemoveAll(tempDir)
	if err != nil {
		log.Fatalf("Failed to write temp dir: %v", err)
	}

	mustRestore(tempDir, assets, nil)
	// listFiles(tempDir)

	err, exitCode := run("./dart", append([]string{"protoc_plugin.snapshot"}, os.Args...), tempDir, []string{})
	if err != nil {
		log.Printf("%v", err)
	}
	os.Exit(exitCode)
}

// MustRestore - Restore assets that are bundled inside the binary.
func mustRestore(baseDir string, assets map[string][]byte, mappings map[string]string) {
	for basename, data := range assets {
		if mappings != nil {
			replacement := mappings[basename]
			if replacement != "" {
				basename = replacement
			}
		}
		filename := path.Join(baseDir, basename)
		dirname := path.Dir(filename)
		if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
			log.Fatalf("Failed to create asset dir %s: %v", dirname, err)
		}

		if err := ioutil.WriteFile(filename, data, os.ModePerm); err != nil {
			log.Fatalf("Failed to write asset %s: %v", filename, err)
		}
	}
}

// listFiles - convenience debugging function to log the files under a given dir
func listFiles(dir string) error {
	log.Println("Listing files under " + dir)
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("%v\n", err)
			return err
		}
		log.Println(path)
		return nil
	})
}

// run a command
func run(entrypoint string, args []string, dir string, env []string) (error, int) {
	cmd := exec.Command(entrypoint, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	cmd.Dir = dir
	err := cmd.Run()

	var exitCode int
	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			log.Printf("Could not get exit code for failed program: %v, %v", entrypoint, args)
			exitCode = -1
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	return err, exitCode
}
