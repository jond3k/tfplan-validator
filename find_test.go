package tfplan_validator

import (
	"os"
	"path"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestFindWorkspaces(t *testing.T) {
	cases := []struct {
		name     string
		in       []string
		expected []string
		errStr   string
	}{
		{
			name:     "empty",
			in:       []string{},
			expected: []string{},
		},
		{
			name:     "no-match",
			in:       []string{"no-match"},
			expected: nil,
			errStr:   "file does not exist",
		},
		{
			name:     "ignore-empty-line",
			in:       []string{""},
			expected: []string{},
		},
		{
			name:     "match-exact",
			in:       []string{"main_only/main.tf"},
			expected: []string{"main_only"},
		},
		{
			name:     "match-negation",
			in:       []string{"main_only/main.tf", "!main_only/main.tf"},
			expected: []string{},
		},
		{
			name:     "match-duplicate",
			in:       []string{"main_and_lock/main.tf", "main_and_lock/.terraform.lock.hcl"},
			expected: []string{"main_and_lock"},
		},
		{
			name:     "default-globs",
			in:       DefaultGlobs,
			expected: []string{"main_and_lock", "main_only"},
		},
	}

	basepath := path.Join(".", "test", "find")

	// Change working dir for glob
	oldwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	} else if err := os.Chdir(basepath); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(oldwd)

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := findWorkspaces(tc.in)
			errStr := makeErrStr(err)
			if !reflect.DeepEqual(tc.expected, actual) || tc.errStr != errStr {
				t.Fatalf("expected:\n\n%s\ngot:\n\n%s\n\nexpected err: %s\n\ngot err: %s\n", spew.Sdump(tc.expected), spew.Sdump(actual), tc.errStr, errStr)
			}
		})
	}

}
