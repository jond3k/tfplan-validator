package tfplan_validator

import (
	tfjson "github.com/hashicorp/terraform-json"
)

// FilterResults describes whether a plan is allowed or not
type FilterResults struct {
	Errors     map[Address]Action
	Changes    map[Address]Action
	PlanFilter *PlanFilter
}

func NewFilterResults(filter *PlanFilter) *FilterResults {
	return &FilterResults{
		Errors:     map[Address]Action{},
		Changes:    map[Address]Action{},
		PlanFilter: filter,
	}
}

// HasChanges returns true if there's at least 1 change
func (fr *FilterResults) HasChanges() bool {
	return len(fr.Changes) > 0
}

// HasErrors returns true if there's at least 1 error
func (fr *FilterResults) HasErrors() bool {
	return len(fr.Errors) > 0
}

// CheckPlan
func CheckPlan(filter *PlanFilter, plan *tfjson.Plan) (*FilterResults, error) {
	results := NewFilterResults(filter)

	for _, change := range plan.ResourceChanges {
		address := Address(change.Address)
		if relevant, err := IsRelevant(change); err != nil {
			return nil, err
		} else if !relevant {
			continue
		}

		proposed := ConvertAction(&change.Change.Actions)

		if filter.HasAction(address, proposed) {
			results.Changes[address] = proposed
		} else {
			results.Errors[address] = proposed
		}
	}

	return results, nil
}
