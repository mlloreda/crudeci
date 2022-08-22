package main_test

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	crudeci "github.com/mlloreda/crudeci"
)

func TestCloneRepo(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name          string
		repo          string
		tmpDir        string
		tmpDirPattern string
		checkFor      string
		expError      bool
	}{
		{
			name:          "success case",
			tmpDirPattern: "test",
			repo:          "https://github.com/mlloreda/dotfiles",
			checkFor:      "dot_vimrc",
			expError:      false,
		},
		{
			name:          "failure case",
			tmpDirPattern: "test",
			repo:          "https://github.com/mlloreda/dotfiles",
			checkFor:      "invalid_file",
			expError:      true,
		},
		// Disabling. Git protocol urls will only succeed locally with
		// appropriate private key.
		// {
		// 	name:          "failure case (git)",
		// 	tmpDirPattern: "test",
		// 	repo:          "git@github.com:mlloreda/invalid",
		// 	expError:      true,
		// },
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			wd := testCreateTempDir(t, tc.tmpDir, tc.tmpDirPattern)

			targetDir := path.Join(wd, path.Base(tc.repo))
			err := crudeci.CloneRepo(tc.repo, targetDir)
			if err != nil {
				t.Fatalf("error: %v", err)
			}

			_, err = crudeci.Exists(path.Join(targetDir, tc.checkFor))
			if (tc.expError == false && err != nil) ||
				(tc.expError == true && err == nil) {

				t.Fatalf("error: %v", err)
			}
		})
	}
}

func testCreateTempDir(t *testing.T, dir, pattern string) string {
	t.Helper()

	wd, err := ioutil.TempDir(dir, pattern)
	if err != nil {
		log.Fatalf("%s", err)
	}
	defer os.RemoveAll(wd)

	return wd
}
