package tfplan_validator

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	tfjson "github.com/hashicorp/terraform-json"
)

// # local_file.foo will be created
// + resource "local_file" "foo" {
// 		+ content              = "foo!"
// 		+ directory_permission = "0777"
// 		+ file_permission      = "0777"
// 		+ filename             = "./foo.bar"
// 		+ id                   = (known after apply)
// 	}

// Plan: 1 to add, 0 to change, 0 to destroy.

// If we like the results we can create a validator that will only accept plans with this create operation. The validator currently only accepts plans in json format.

// > terraform show -json ./plan > ./plan.json
// > tfplan-validator create ../rules.json ./plan.json

// Created rules file ../rules.json that allows Terraform to perform the following actions:

// - local_file.foo can be created

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
			in:   readPlansP([]string{"fixtures/create/plan.json"}),
			expected: &PlanFilter{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"local_file.foo": {ActionCreate},
				},
			},
		},
		{
			name: "create-delete",
			in:   readPlansP([]string{"fixtures/create-delete/plan.json"}),
			expected: &PlanFilter{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"local_file.foo": {ActionCreateBeforeDestroy},
				},
			},
		},
		{
			name: "delete",
			in:   readPlansP([]string{"fixtures/delete/plan.json"}),
			expected: &PlanFilter{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"local_file.foo": {ActionDelete},
				},
			},
		},
		{
			name: "delete-create",
			in:   readPlansP([]string{"fixtures/delete-create/plan.json"}),
			expected: &PlanFilter{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"local_file.foo": {ActionDestroyBeforeCreate},
				},
			},
		},
		{
			name: "update",
			in:   readPlansP([]string{"fixtures/update/plan.json"}),
			expected: &PlanFilter{
				FormatVersion: CurrentFormatVersion,
				AllowedActions: map[Address][]Action{
					"google_project_iam_policy.project": {ActionUpdate},
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
