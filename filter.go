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
	if data, err := json.MarshalIndent(pf, "", "  "); err != nil {
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

// relevantActions are the only types we care about
var relevantActions = map[Action]bool{
	ActionCreate:              true,
	ActionUpdate:              true,
	ActionDelete:              true,
	ActionDestroyBeforeCreate: true,
	ActionCreateBeforeDestroy: true,
}

// IsRelevant returns true if the change is something we care about
func IsRelevant(rc *tfjson.ResourceChange) (bool, error) {
	if rc.Mode != tfjson.ManagedResourceMode {
		return false, nil
	}
	if action := ConvertAction(&rc.Change.Actions); action != ActionInvalid {
		return relevantActions[action], nil
	}
	return false, fmt.Errorf("unrecognized change in plan: %v", rc)
}

// FilterForPlan creates a filter that accepts everything in the specified plan
func NewFilterFromPlan(plan *tfjson.Plan) (*PlanFilter, error) {
	allowed := map[Address][]Action{}

	for _, rc := range plan.ResourceChanges {
		if b, err := IsRelevant(rc); err != nil {
			return nil, err
		} else if !b {
			continue
		}

		address := Address(rc.Address)

		if current := allowed[address]; current != nil {
			return nil, fmt.Errorf("duplicate address in plan: %s", current)
		}

		allowed[address] = []Action{ConvertAction(&rc.Change.Actions)}
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
