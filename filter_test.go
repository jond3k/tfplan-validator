package tfplan_validator

import (
	"path"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	tfjson "github.com/hashicorp/terraform-json"
)

func planPath(typ string) string {
	return path.Join("examples", "fixtures", typ, "plan.json")
}

func filterPath(typ string) string {
	return path.Join("examples", "fixtures", typ, "filter.json")
}

func otherPath(name string) string {
	return path.Join("examples", "fixtures", "itest", name)
}

func readPlansP(paths []string) []*tfjson.Plan {
	if plans, err := ReadPlans(paths); err != nil {
		panic(err)
	} else {
		return plans
	}
}

func TestNewFilterFromPlans(t *testing.T) {
	cases := []struct {
		name     string
		in       []*tfjson.Plan
		expected *PlanFilter
		err      error
	}{
		{
			name: "empty",
			in:   []*tfjson.Plan{{}},
			expected: &PlanFilter{
				FormatVersion:  CurrentFormatVersion,
				AllowedActions: map[Address][]Action{},
			},
		},
		{
			name: "create",
			in:   readPlansP([]string{planPath("create")}),
			expected: &PlanFilter{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"local_file.foo": {ActionCreate},
				},
			},
		},
		{
			name: "create-delete",
			in:   readPlansP([]string{planPath("create-delete")}),
			expected: &PlanFilter{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"local_file.foo": {ActionCreateBeforeDestroy},
				},
			},
		},
		{
			name: "delete",
			in:   readPlansP([]string{planPath("delete")}),
			expected: &PlanFilter{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"local_file.foo": {ActionDelete},
				},
			},
		},
		{
			name: "delete-create",
			in:   readPlansP([]string{planPath("delete-create")}),
			expected: &PlanFilter{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"local_file.foo": {ActionDestroyBeforeCreate},
				},
			},
		},
		{
			name: "update",
			in:   readPlansP([]string{planPath("update")}),
			expected: &PlanFilter{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"google_project_iam_policy.project": {ActionUpdate},
				},
			},
		},
		{
			name: "ignore data",
			in: []*tfjson.Plan{
				{
					ResourceChanges: []*tfjson.ResourceChange{
						{
							Mode:    tfjson.DataResourceMode,
							Address: "a.b.c",
						},
					},
				},
			},
			expected: &PlanFilter{
				FormatVersion:  CurrentFormatVersion,
				AllowedActions: map[Address][]Action{},
			},
		},
		{
			name: "create and create-delete are compatible",
			in: []*tfjson.Plan{
				{
					ResourceChanges: []*tfjson.ResourceChange{
						{
							Mode:    tfjson.ManagedResourceMode,
							Address: "a.b.c",
							Change: &tfjson.Change{
								Actions: tfjson.Actions{tfjson.ActionCreate},
							},
						},
					},
				},
				{
					ResourceChanges: []*tfjson.ResourceChange{
						{
							Mode:    tfjson.ManagedResourceMode,
							Address: "a.b.c",
							Change: &tfjson.Change{
								Actions: tfjson.Actions{tfjson.ActionCreate, tfjson.ActionDelete},
							},
						},
					},
				},
			},
			expected: &PlanFilter{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"a.b.c": {ActionCreate, ActionCreateBeforeDestroy},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if actual, err := NewFilterFromPlans(tc.in); err != nil {
				t.Fatal(err)
			} else if !reflect.DeepEqual(tc.expected, actual) {
				t.Fatalf("expected:\n\n%s\ngot:\n\n%s\n", spew.Sdump(tc.expected), spew.Sdump(actual))
			}
		})
	}
}

func TestReadPlanFilters(t *testing.T) {
	cases := []struct {
		name     string
		in       []string
		expected []*PlanFilter
		errStr   string
	}{
		{
			name: "load two files",
			in:   []string{filterPath("create"), filterPath("update")},
			expected: []*PlanFilter{
				{
					FormatVersion: CurrentFormatVersion,
					AllowedActions: map[Address][]Action{
						"local_file.foo": {ActionCreate},
					},
				},
				{
					FormatVersion: CurrentFormatVersion,
					AllowedActions: map[Address][]Action{
						"google_project_iam_policy.project": {ActionUpdate},
					},
				},
			},
		},
		{
			name:   "one file is missing",
			in:     []string{filterPath("create"), filterPath("missing")},
			errStr: "open " + filterPath("missing") + ": no such file or directory",
		},
		{
			name:   "one file is missing",
			in:     []string{filterPath("create"), otherPath("unparseable.json")},
			errStr: otherPath("unparseable.json") + ": unexpected end of JSON input",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := ReadPlanFilters(tc.in)
			errStr := makeErrStr(err)
			if !reflect.DeepEqual(tc.expected, actual) || tc.errStr != errStr {
				t.Fatalf("expected:\n\n%s\ngot:\n\n%s\n\nexpected err:%s\n\ngot err: %s\n", spew.Sdump(tc.expected), spew.Sdump(actual), tc.errStr, errStr)
			}
		})
	}
}
