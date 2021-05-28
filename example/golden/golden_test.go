package golden

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/stackb/rules_proto/pkg/goldentest"
)

func TestProtoc(t *testing.T) {
	// listFiles(".")
	goldentest.FromDir("example/golden").Run(t, "gazelle")
}

// listFiles - convenience debugging function to log the files under a given dir
func listFiles(dir string) error {
	return filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("%v\n", err)
			return err
		}
		if info.Mode()&os.ModeSymlink > 0 {
			link, err := os.Readlink(path)
			if err != nil {
				return err
			}
			log.Printf("%s -> %s", path, link)
			return nil
		}

		log.Println(path)
		return nil
	})
}
