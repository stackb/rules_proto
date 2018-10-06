package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/golang/protobuf/proto"
	compiler_plugin "google.golang.org/protobuf/compiler_plugin"
)

func main() {
	tempDir, err := ioutil.TempDir("", "protoc-gen-swift-")
	defer os.RemoveAll(tempDir)
	if err != nil {
		log.Fatalf("Failed to write temp dir: %v", err)
	}
	files := MustRestore(tempDir, assets, nil)

	var compiler string
	for _, file := range files {
		if strings.HasSuffix(file, "ProtoCompilerPlugin") {
			compiler = file
			break
		}
	}

	if compiler == "" {
		ListFiles(tempDir)
		log.Fatalf("Failed to locate the compiler plugin!")
	}

	var request compiler_plugin.CodeGeneratorRequest
	bytes, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatalf("Failed to read stdin: %v", err)
	}

	proto.Unmarshal(bytes, &request)
	log.Printf("request: %+v", request)

	requestFile := path.Join(tempDir, "request.proto")
	err = ioutil.WriteFile("request.proto", bytes, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to write request proto: %v", err)
	}

	//err, exitCode := Run(compiler, os.Args, tempDir, nil)
	err, exitCode := Run(compiler, []string{requestFile}, tempDir, nil)

	if err != nil {
		log.Printf("%v", err)
	}

	os.Exit(exitCode)
}

// Run a command
func Run(entrypoint string, args []string, dir string, env []string) (error, int) {
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

// MustRestore - Restore assets.
func MustRestore(baseDir string, assets map[string][]byte, mappings map[string]string) []string {
	files := make([]string, 0)

	// unpack variable is provided by the go_embed data and is a
	// map[string][]byte such as {"/usr/share/games/fortune/literature.dat":
	// bytes... }
	for basename, data := range assets {
		if mappings != nil {
			replacement := mappings[basename]
			if replacement != "" {
				basename = replacement
			}
		}
		// If not a tarball, write file directly
		filename := path.Join(baseDir, basename)
		dirname := path.Dir(filename)
		//log.Printf("file %s, dir %s, rel %d, abs %s, absdir: %s", file, dir, rel, abs, absdir)
		if err := os.MkdirAll(dirname, os.ModePerm); err != nil {
			log.Fatalf("Failed to create asset dir %s: %v", dirname, err)
		}

		if err := ioutil.WriteFile(filename, data, os.ModePerm); err != nil {
			log.Fatalf("Failed to write asset %s: %v", filename, err)
		}
		//log.Printf("Restored %s", filename)
		files = append(files, filename)
	}

	return files
	//log.Printf("Assets restored to %s", baseDir)
}

// ListFiles - convenience debugging function to log the files under a given dir
func ListFiles(dir string) error {
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
