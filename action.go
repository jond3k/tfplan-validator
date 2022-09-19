package tfplan_validator

import (
	tfjson "github.com/hashicorp/terraform-json"
)

// Action is an enum that corresponds to a valid tfjson plan array e.g. ["delete", "create"] is ActionDeleteCreate
type Action string

const (
	ActionInvalid             Action = "invalid"
	ActionNoOp                Action = "no-op"
	ActionRead                Action = "read"
	ActionCreate              Action = "create"
	ActionUpdate              Action = "update"
	ActionDelete              Action = "delete"
	ActionDestroyBeforeCreate Action = "delete-create"
	ActionCreateBeforeDestroy Action = "create-delete"
)

// IsEqual to a TF plan action?
func (a Action) IsEqual(actions *tfjson.Actions) bool {
	return ConvertAction(actions) == a
}

// ConvertAction from the tfjson form to one we can more easily work with
func ConvertAction(actions *tfjson.Actions) Action {
	if actions.NoOp() {
		return ActionNoOp
	} else if actions.Read() {
		return ActionRead
	} else if actions.Create() {
		return ActionCreate
	} else if actions.Update() {
		return ActionUpdate
	} else if actions.Delete() {
		return ActionDelete
	} else if actions.DestroyBeforeCreate() {
		return ActionDestroyBeforeCreate
	} else if actions.CreateBeforeDestroy() {
		return ActionCreateBeforeDestroy
	}
	return ActionInvalid
}