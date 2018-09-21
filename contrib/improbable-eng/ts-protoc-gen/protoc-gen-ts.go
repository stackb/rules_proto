package main

import (
	"archive/tar"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"syscall"
)

func main() {
	tempDir, err := ioutil.TempDir("", "protoc-gen-ts-")
	defer os.RemoveAll(tempDir)
	if err != nil {
		log.Fatalf("Failed to write temp dir: %v", err)
	}
	MustRestore(tempDir, assets, nil)
	//ListFiles(tempDir)

	err, exitCode := Run("external/node/bin/node", append([]string{"./ts_protoc_gen/src/index.js"}, os.Args...), tempDir, []string{
		fmt.Sprintf("NODE_PATH=%s:%s/deps/node_modules", tempDir, tempDir),
	})
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
func MustRestore(baseDir string, assets map[string][]byte, mappings map[string]string) {
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
		if strings.HasSuffix(basename, ".tar") {
			// Untar any tarballs directly to baseDir
			if err := Untar(bytes.NewReader(data), baseDir); err != nil {
				log.Fatalf("Failed to untar asset tarball %s: %v", basename, err)
			}
			//log.Printf("Untarred %s", basename)
			continue
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
	}

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

// Untar takes a destination path and a reader; a tar reader loops over the
// tarfile creating the file structure at 'dst' along the way, and writing any
// files
func Untar(r io.Reader, dst string) error {

	// gzr, err := gzip.NewReader(r)
	// defer gzr.Close()
	// if err != nil {
	// 	return err
	// }

	tr := tar.NewReader(r)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}
			defer f.Close()

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}
			// Even though we defer, try closing now lest we get a 'too many
			// open files' error
			f.Close()
		}
	}
}
