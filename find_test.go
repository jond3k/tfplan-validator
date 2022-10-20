package tfplan_validator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func TestFindWorkspaces(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name       string
		files      []string
		searchDirs []string
		globs      []string
		expected   []string
		errStr     string
	}{
		{
			name:       "error-if-no-globs",
			files:      []string{},
			searchDirs: []string{"."},
			globs:      []string{},
			expected:   nil,
			errStr:     "expected at least one glob",
		},
		{
			name:       "error-if-no-search-dirs",
			files:      []string{},
			searchDirs: []string{},
			globs:      []string{"glob"},
			expected:   nil,
			errStr:     "expected at least one searchDir",
		},
		{
			name:       "nothing-in-empty-folder",
			files:      []string{},
			searchDirs: []string{"."},
			globs:      []string{"**"},
			expected:   []string{},
		},
		{
			name:       "everything-in-folder",
			files:      []string{"a/main.tf", "b/main.tf", "c/main.tf"},
			searchDirs: []string{"."},
			globs:      []string{"**/main.tf"},
			expected:   []string{"a", "b", "c"},
		},
		{
			name:       "ignore-empty-line-glob",
			files:      []string{"a/main.tf", "b/main.tf", "c/main.tf"},
			searchDirs: []string{"."},
			globs:      []string{"", "**/main.tf"},
			expected:   []string{"a", "b", "c"},
		},
		{
			name:       "multiple-folders",
			files:      []string{"a/a/main.tf", "b/b/main.tf", "c/main.tf"},
			searchDirs: []string{"a", "b"},
			globs:      []string{"**/main.tf"},
			expected:   []string{"a/a", "b/b"},
		},
		{
			name:       "default-matches-main-tf",
			files:      []string{"a/main.tf", "b/ignore.txt"},
			searchDirs: []string{"."},
			globs:      DefaultGlobs,
			expected:   []string{"a"},
		},
		{
			name:       "default-matches-tf-lock",
			files:      []string{"a/.terraform.lock.hcl", "b/ignore.txt"},
			searchDirs: []string{"."},
			globs:      DefaultGlobs,
			expected:   []string{"a"},
		},
		{
			name:       "default-matches-main-and-lock-no-duplicate",
			files:      []string{"a/main.tf", "a/.terraform.lock.hcl", "b/ignore.txt"},
			searchDirs: []string{"."},
			globs:      DefaultGlobs,
			expected:   []string{"a"},
		},
		{
			name:       "default-ignores-modules",
			files:      []string{"modules/a/main.tf", "b/main.tf", "b/ignore.txt"},
			searchDirs: []string{"."},
			globs:      DefaultGlobs,
			expected:   []string{"b"},
		},
		{
			name:       "default-ignores-terragrunt-cache",
			files:      []string{"a/main.tf", "a/.terragrunt-cache/main.tf", "b/ignore.txt"},
			searchDirs: []string{"."},
			globs:      DefaultGlobs,
			expected:   []string{"a"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			basepath := path.Join(wd, "test-results", fmt.Sprintf("find-test-%d", time.Now().Unix()), tc.name)

			// create test data
			if err := os.MkdirAll(basepath, 0700); err != nil {
				t.Fatal(err)
			}
			for idx, fname := range tc.files {
				fname = path.Join(basepath, fname)
				dname := filepath.Dir(fname)
				if err := os.MkdirAll(dname, 0700); err != nil {
					t.Fatal(err)
				} else if err := ioutil.WriteFile(fname, []byte{}, 0700); err != nil {
					t.Fatal(err)
				}
				tc.files[idx] = fname
			}

			// normalize search dirs
			for idx, searchDir := range tc.searchDirs {
				tc.searchDirs[idx] = path.Join(basepath, searchDir)
			}

			// normalize expected dirs
			for idx, exp := range tc.expected {
				tc.expected[idx] = path.Join(basepath, exp)
			}

			actual, err := FindWorkspaces(tc.searchDirs, tc.globs)
			errStr := makeErrStr(err)
			if !reflect.DeepEqual(tc.expected, actual) || tc.errStr != errStr {
				t.Fatalf("expected:\n\n%s\ngot:\n\n%s\n\nexpected err: %s\n\ngot err: %s\n", spew.Sdump(tc.expected), spew.Sdump(actual), tc.errStr, errStr)
			}
		})
	}

}
