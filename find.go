package tfplan_validator

import (
	"errors"
	"os"
	"path/filepath"
	"sort"

	"github.com/mattn/go-zglob"
)

var DefaultGlobs = []string{
	"**/main.tf",
	"**/.terraform.lock.hcl",
	"!**/modules/**/main.tf",
	"!**/modules/**/.terraform.lock.hcl",
	"!**/.terragrunt-cache/**",
}

// FindWorkspaces iterates the current working directory and finds candidate workspaces
// it takes a series of globs which support double stars for recursion and ! for negation
func FindWorkspaces(searchDirs []string, globs []string) ([]string, error) {
	paths := map[string]bool{}

	if len(searchDirs) < 1 {
		return nil, errors.New("expected at least one searchDir")
	} else if len(globs) < 1 {
		return nil, errors.New("expected at least one glob")
	}

	oldwd, err := os.Getwd()

	if err != nil {
		return nil, err
	}

	for i, dir := range searchDirs {
		if searchDirs[i], err = filepath.Abs(dir); err != nil {
			return nil, err
		}
	}

	for _, dir := range searchDirs {
		for _, glob := range globs {
			if len(glob) < 1 {
				continue
			}
			negate := glob[0] == '!'
			if negate {
				glob = glob[1:]
			}

			if err := os.Chdir(dir); err != nil {
				return nil, err
			}
			defer os.Chdir(oldwd)

			files, err := zglob.Glob(glob)
			if err != nil {
				return nil, err
			}
			for _, file := range files {
				if absdir, err := filepath.Abs(filepath.Dir(file)); err != nil {
					return nil, err
				} else {
					paths[absdir] = !negate
				}
			}
		}
	}

	results := []string{}

	for k, v := range paths {
		if v {
			results = append(results, k)
		}
	}

	sort.Slice(results, func(i, j int) bool { return results[i] < results[j] })
	return results, nil
}
