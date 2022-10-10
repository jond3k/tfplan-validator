package tfplan_validator

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func TestFindWorkspaces(t *testing.T) {
	files := []string{
		// Test cases in depth 1
		"empty/.keep",
		"main_only/main.tf",
		"lock_only/.terraform.lock.hcl",
		"modules/module1/main.tf",
		"both/main.tf",
		"both/.terraforn.lock.hcl",
		// Test cases at depth 3
		"a/b/c/empty/.keep",
		"a/b/c/main_only/main.tf",
		"a/b/c/lock_only/.terraform.lock.hcl",
		"a/b/c/modules/module1/main.tf",
		"a/b/c/both/main.tf",
		"a/b/c/both/.terraforn.lock.hcl",
	}
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
		// !match
		// match
		// de-dupe
		// sequential
	}

	basepath := path.Join(".", "test-results", "find", fmt.Sprint(time.Now().Unix()))

	for _, file := range files {
		filepath := path.Join(basepath, file)
		if err := os.MkdirAll(path.Dir(filepath), 0700); err != nil {
			t.Fatal(err)
		} else if err := ioutil.WriteFile(filepath, []byte{}, 0700); err != nil {
			t.Fatal(err)
		}
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := findWorkspaces(tc.in)
			errStr := makeErrStr(err)
			if !reflect.DeepEqual(tc.expected, actual) || tc.errStr != errStr {
				t.Fatalf("expected:\n\n%s\ngot:\n\n%s\n\nexpected err:%s\n\ngot err: %s\n", spew.Sdump(tc.expected), spew.Sdump(actual), tc.errStr, errStr)
			}
		})
	}

}
