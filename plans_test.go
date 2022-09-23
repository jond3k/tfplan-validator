package tfplan_validator

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	tfjson "github.com/hashicorp/terraform-json"
)

func TestReadPlans(t *testing.T) {
	cases := []struct {
		name     string
		in       []string
		expected []*tfjson.Plan
		errStr   string
	}{
		{
			name:     "load two files",
			in:       []string{planPath("create"), planPath("update")},
			expected: readPlansP([]string{planPath("create"), planPath("update")}),
		},
		{
			name:   "one file is missing",
			in:     []string{planPath("create"), otherPath("missing.json")},
			errStr: "open " + otherPath("missing.json") + ": no such file or directory",
		},
		{
			name:   "one file is invalid json",
			in:     []string{planPath("create"), otherPath("unparseable.json")},
			errStr: otherPath("unparseable.json") + ": unexpected end of JSON input",
		},
		{
			name:   "one file has invalid content",
			in:     []string{planPath("create"), otherPath("plan-missing-version.json")},
			errStr: otherPath("plan-missing-version.json") + ": unexpected plan input, format version is missing",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ReadPlans(tc.in)
			errStr := makeErrStr(err)
			if !reflect.DeepEqual(tc.expected, actual) || tc.errStr != errStr {
				t.Fatalf("expected:\n\n%s\ngot:\n\n%s\n\nexpected err:%s\n\ngot err: %s\n", spew.Sdump(tc.expected), spew.Sdump(actual), tc.errStr, errStr)
			}
		})
	}
}
