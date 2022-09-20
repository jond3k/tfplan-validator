package tfplan_validator

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	tfjson "github.com/hashicorp/terraform-json"
)

// CurrentFormatVersion of the rules file
const CurrentFormatVersion = "0.1"

// Address of a resource in a plan
type Address string

// PlanFilter defines what is allowed in a plan
type PlanFilter struct {
	// The version of the PlanFilter format
	FormatVersion string `json:"format_version,omitempty"`

	// The actions we can perform on a resource
	AllowedActions map[Address][]Action `json:"allowed_actions"`
}

// ReadPlanFilter from a path
func ReadPlanFilter(path string) (*PlanFilter, error) {
	if data, err := ioutil.ReadFile(path); err != nil {
		return nil, err
	} else {
		return ParsePlanFilter(data)
	}
}

// WritePlanFilter to a path
func (pf *PlanFilter) WriteJSON(path string) error {
	if data, err := json.MarshalIndent(pf, "", " "); err != nil {
		return err
	} else if err := ioutil.WriteFile(path, data, 0644); err != nil {
		return err
	}
	return nil
}

// ParsePlanFilter from JSON
func ParsePlanFilter(data []byte) (*PlanFilter, error) {
	var f PlanFilter
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// IsRelevant returns true if the change is something we care about
func IsRelevant(rc *tfjson.ResourceChange) (bool, error) {
	if rc.Mode != tfjson.ManagedResourceMode {
		return false, nil
	}
	switch ConvertAction(&rc.Change.Actions) {
	case ActionNoOp:
		return false, nil
	case ActionRead:
		return false, nil
	case ActionCreate:
		return true, nil
	case ActionUpdate:
		return true, nil
	case ActionDelete:
		return true, nil
	case ActionDestroyBeforeCreate:
		return true, nil
	case ActionCreateBeforeDestroy:
		return true, nil
	}
	return false, fmt.Errorf("unrecognized change in plan: %v", rc)
}

// FilterForPlan creates a filter that accepts everything in the specified plan
func NewFilterFromPlan(plan *tfjson.Plan) (*PlanFilter, error) {
	allowed := map[Address][]Action{}

	for _, rc := range plan.ResourceChanges {
		if rc.Mode != tfjson.ManagedResourceMode {
			continue
		}

		if b, err := IsRelevant(rc); err != nil {
			return nil, err
		} else if !b {
			continue
		}

		address := Address(rc.Address)
		var current []Action
		if current = allowed[Address(address)]; current == nil {
			current = []Action{}
		}

		action := ConvertAction(&rc.Change.Actions)

		for _, other := range current {
			if !AreCompatible(action, other) {
				return nil, fmt.Errorf("contradictory actions: %s has %s and %s", address, action, other)
			}
		}

		allowed[address] = append(current, action)
	}

	return &PlanFilter{
		FormatVersion:  CurrentFormatVersion,
		AllowedActions: allowed,
	}, nil
}

// NewFilterFromPlans creates a filter that accepts everything in a list of plans
func NewFilterFromPlans(plans []*tfjson.Plan) (*PlanFilter, error) {
	filters := make([]*PlanFilter, len(plans))
	for i, plan := range plans {
		if filter, err := NewFilterFromPlan(plan); err != nil {
			return nil, err
		} else {
			filters[i] = filter
		}
	}
	return MergePlanFilters(filters)
}

// NewFilterFromPlanPaths creates filter from a sequence of paths
func NewFilterFromPlanPaths(paths []string) (*PlanFilter, error) {
	if plans, err := ReadPlans(paths); err != nil {
		return nil, err
	} else {
		return NewFilterFromPlans(plans)
	}
}

// MergePlanFilters combines filters from multiple plans into one
func MergePlanFilters(filters []*PlanFilter) (*PlanFilter, error) {
	// TODO: currently only supporting one
	return filters[0], nil
}
