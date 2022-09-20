package tfplan_validator

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestMergePlanFilters(t *testing.T) {
	cases := []struct {
		name     string
		in       []*PlanFilter
		expected *PlanFilter
		err      error
	}{
		{
			name: "empty",
			in: []*PlanFilter{{
				FormatVersion:  CurrentFormatVersion,
				AllowedActions: map[Address][]Action{},
			}, {
				FormatVersion:  CurrentFormatVersion,
				AllowedActions: map[Address][]Action{},
			}},
			expected: &PlanFilter{
				FormatVersion:  CurrentFormatVersion,
				AllowedActions: map[Address][]Action{},
			},
		},
		{
			name: "different keys",
			in: []*PlanFilter{{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"a.b.c": {ActionCreate},
				},
			}, {
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"d.e.f": {ActionUpdate},
				},
			}},
			expected: &PlanFilter{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"a.b.c": {ActionCreate},
					"d.e.f": {ActionUpdate},
				},
			},
		},
		{
			name: "duplicates",
			in: []*PlanFilter{{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"a.b.c": {ActionCreate},
					"d.e.f": {ActionUpdate},
				},
			}, {
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"a.b.c": {ActionCreate},
				},
			}},
			expected: &PlanFilter{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"a.b.c": {ActionCreate},
					"d.e.f": {ActionUpdate},
				},
			},
		},
		{
			name: "combine if compatible",
			in: []*PlanFilter{{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"a.b.c": {ActionCreate},
				},
			}, {
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"a.b.c": {ActionUpdate},
				},
			}},
			expected: &PlanFilter{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"a.b.c": {ActionCreate, ActionUpdate},
				},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if actual, err := MergePlanFilters(tc.in); err != nil {
				t.Fatal(err)
			} else if !reflect.DeepEqual(tc.expected, actual) {
				t.Fatalf("expected:\n\n%s\ngot:\n\n%s\n", spew.Sdump(tc.expected), spew.Sdump(actual))
			}
		})
	}
}
