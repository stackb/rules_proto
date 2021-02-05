package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProtocopierErrors(t *testing.T) {
	cases := []struct {
		d      string
		cfg    *Config
		inputs map[string]bool
		err    error
	}{
		{
			d:   "complains about empty inputs",
			cfg: &Config{},
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			//
			// Setup the source directory with mock files
			//
			srcDir, err := ioutil.TempDir("src", "")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(srcDir)

			for rel := range tc.inputs {
				abs := filepath.Join(srcDir, rel)
				if err := ioutil.WriteFile(abs, nil, os.ModePerm); err != nil {
					t.Fatal(err)
				}
			}

			//
			// Setup the destination dir where files will be created
			//
			dstDir, err := ioutil.TempDir("dst", "")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(srcDir)

			tc.cfg.BuildWorkspaceDirectory = dstDir

			if err := run(tc.cfg); err != nil {

			}
			require.Equal(t, tc.err, err)
		})
	}
}
