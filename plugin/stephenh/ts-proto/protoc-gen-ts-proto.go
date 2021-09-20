package main

// protoc-gen-ts-proto.go implements a binary that coordinates calling the
// nodejs plugin entrypoint.  nodejs_binary does not work as a 'tool'.

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func main() {
	// environment is stamped into the binary as an .env file; first we need to parse it
	//
	env := make(map[string]string)

	for _, data := range assets {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			env[parts[0]] = parts[1]
		}

	}

	if env["NODE_BIN"] == "" {
		log.Fatal("NODE_BIN is required")
	}
	if env["NPM_WORKSPACE_BUILD_FILE"] == "" {
		log.Fatal("NPM_WORKSPACE_BUILD_FILE is required")
	}

	npmWorkspace := filepath.Dir(env["NPM_WORKSPACE_BUILD_FILE"])
	entrypoint := filepath.Join(".", npmWorkspace, "node_modules", "ts-proto/build/plugin.js")

	exitCode, err := run(env["NODE_BIN"], []string{
		"--eval",
		fmt.Sprintf(`require("./%s")`, entrypoint),
	}, ".", nil)
	if err != nil {
		log.Println(err)
	}
	os.Exit(exitCode)
}

func run(entrypoint string, args []string, dir string, env []string) (int, error) {
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
			// This will happen (in OSX) if `name` is not available in $PATH, in
			// this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and
			// format err to string and set to stderr
			log.Printf("Could not get exit code for failed program: %v, %v", entrypoint, args)
			exitCode = -1
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}

	return exitCode, err
}
