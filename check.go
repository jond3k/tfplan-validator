package tfplan_validator

import (
	tfjson "github.com/hashicorp/terraform-json"
)

// FilterResults describes whether a plan is allowed or not
type FilterResults struct {
	Errors  map[Address]Action
	Changes map[Address]Action
}

// HasErrors returns true if there's at least 1 error
func (fr *FilterResults) HasErrors() bool {
	return len(fr.Errors) > 0
}

// CheckPlan
func CheckPlan(filter *PlanFilter, plan *tfjson.Plan) (*FilterResults, error) {
	var results FilterResults

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

	return &results, nil
}
