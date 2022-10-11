package tfplan_validator

import (
	"path"
	"sort"

	"github.com/mattn/go-zglob"
)

var DefaultGlobs = []string{
	"**/main.tf",
	"**/.terraform.lock.hcl",
	"!**/modules/**/main.tf",
	"!**/modules/**/.terraform.lock.hcl",
}

// FindWorkspaces iterates the current working directory and finds candidate workspaces
// it takes a series of globs which support double stars for recursion and ! for negation
func FindWorkspaces(globs []string) ([]string, error) {
	paths := map[string]bool{}
	for _, glob := range globs {
		if len(glob) < 1 {
			continue
		}
		negate := glob[0] == '!'
		if negate {
			glob = glob[1:]
		}
		files, err := zglob.Glob(glob)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			paths[path.Dir(file)] = !negate
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
