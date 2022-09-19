package tfplan_validator

import (
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
	tfjson "github.com/hashicorp/terraform-json"
)

func TestConvertAction(t *testing.T) {
	cases := []struct {
		name     string
		in       *tfjson.Actions
		expected Action
		err      error
	}{
		{
			name:     "no-op",
			in:       &tfjson.Actions{tfjson.ActionNoop},
			expected: ActionNoOp,
		},
		{
			name:     "read",
			in:       &tfjson.Actions{tfjson.ActionRead},
			expected: ActionRead,
		},
		{
			name:     "create",
			in:       &tfjson.Actions{tfjson.ActionCreate},
			expected: ActionCreate,
		},
		{
			name:     "update",
			in:       &tfjson.Actions{tfjson.ActionUpdate},
			expected: ActionUpdate,
		},
		{
			name:     "delete",
			in:       &tfjson.Actions{tfjson.ActionDelete},
			expected: ActionDelete,
		},
		{
			name:     "delete-create",
			in:       &tfjson.Actions{tfjson.ActionDelete, tfjson.ActionCreate},
			expected: ActionDestroyBeforeCreate,
		},
		{
			name:     "create-delete",
			in:       &tfjson.Actions{tfjson.ActionCreate, tfjson.ActionDelete},
			expected: ActionCreateBeforeDestroy,
		},
		{
			name:     "invalid",
			in:       &tfjson.Actions{tfjson.ActionCreate, tfjson.ActionCreate},
			expected: ActionInvalid,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if actual := ConvertAction(tc.in); !reflect.DeepEqual(tc.expected, actual) {
				t.Fatalf("expected:\n\n%s\ngot:\n\n%s\n", spew.Sdump(tc.expected), spew.Sdump(actual))
			}
		})
	}
}
