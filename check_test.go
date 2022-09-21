package tfplan_validator

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	tfjson "github.com/hashicorp/terraform-json"
)

func makeErrStr(err error) string {
	if err == nil {
		return ""
	}
	return err.Error()
}

func TestHasChanges(t *testing.T) {
	cases := []struct {
		name     string
		in       *FilterResults
		expected bool
	}{
		{
			name: "with changes",
			in: &FilterResults{
				Changes: map[Address]Action{
					"a.b.c": ActionCreate,
				},
			},
			expected: true,
		},
		{
			name:     "without changes",
			in:       &FilterResults{},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := tc.in.HasChanges(); actual != tc.expected {
				t.Fatalf("expected:\n\n%s\ngot:\n\n%s\n\nexpected\n", spew.Sdump(tc.expected), spew.Sdump(actual))
			}
		})
	}
}

func TestHasErrors(t *testing.T) {
	cases := []struct {
		name     string
		in       *FilterResults
		expected bool
	}{
		{
			name: "with errors",
			in: &FilterResults{
				Errors: map[Address]Action{
					"a.b.c": ActionCreate,
				},
			},
			expected: true,
		},
		{
			name:     "without errors",
			in:       &FilterResults{},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := tc.in.HasErrors(); actual != tc.expected {
				t.Fatalf("expected:\n\n%s\ngot:\n\n%s\n\nexpected\n", spew.Sdump(tc.expected), spew.Sdump(actual))
			}
		})
	}
}

func TestCheckPlan(t *testing.T) {
	cases := []struct {
		name     string
		filter   *PlanFilter
		plan     *tfjson.Plan
		expected *FilterResults
		errStr   string
	}{
		{
			name: "no changes",
			filter: &PlanFilter{
				FormatVersion:  CurrentFormatVersion,
				AllowedActions: map[Address][]Action{},
			},
			plan: &tfjson.Plan{},
			expected: NewFilterResults(&PlanFilter{
				FormatVersion:  CurrentFormatVersion,
				AllowedActions: map[Address][]Action{},
			}),
		},
		{
			name: "ignore irrelevant changes",
			filter: &PlanFilter{
				FormatVersion:  CurrentFormatVersion,
				AllowedActions: map[Address][]Action{},
			},
			plan: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Address: "a.b.c",
						Mode:    tfjson.ManagedResourceMode,
						Change: &tfjson.Change{
							Actions: tfjson.Actions{tfjson.ActionNoop},
						},
					},
				},
			},
			expected: NewFilterResults(&PlanFilter{
				FormatVersion:  CurrentFormatVersion,
				AllowedActions: map[Address][]Action{},
			}),
		},
		{
			name: "reject invalid actions",
			filter: &PlanFilter{
				FormatVersion:  CurrentFormatVersion,
				AllowedActions: map[Address][]Action{},
			},
			plan: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Address: "a.b.c",
						Mode:    tfjson.ManagedResourceMode,
						Change: &tfjson.Change{
							Actions: tfjson.Actions{"invalid"},
						},
					},
				},
			},
			errStr: "unrecognized action in plan: [invalid]",
		},
		{
			name: "add changes and errors",
			filter: &PlanFilter{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"a.b.c": {ActionCreate},
					"d.e.f": {ActionCreate},
				},
			},
			plan: &tfjson.Plan{
				ResourceChanges: []*tfjson.ResourceChange{
					{
						Address: "a.b.c",
						Mode:    tfjson.ManagedResourceMode,
						Change: &tfjson.Change{
							Actions: tfjson.Actions{tfjson.ActionCreate},
						},
					},
					{
						Address: "d.e.f",
						Mode:    tfjson.ManagedResourceMode,
						Change: &tfjson.Change{
							Actions: tfjson.Actions{tfjson.ActionDelete},
						},
					},
				},
			},
			expected: &FilterResults{
				Changes: map[Address]Action{
					"a.b.c": ActionCreate,
				},
				Errors: map[Address]Action{
					"d.e.f": ActionDelete,
				},
				PlanFilter: &PlanFilter{
					FormatVersion: CurrentFormatVersion,
					AllowedActions: map[Address][]Action{
						"a.b.c": {ActionCreate},
						"d.e.f": {ActionCreate},
					},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual, err := CheckPlan(tc.filter, tc.plan)
			errStr := makeErrStr(err)
			if !reflect.DeepEqual(tc.expected, actual) || tc.errStr != errStr {
				t.Fatalf("expected:\n\n%s\ngot:\n\n%s\n\nexpected err:%s\n\ngot err: %s\n", spew.Sdump(tc.expected), spew.Sdump(actual), tc.errStr, errStr)
			}
		})
	}
}
