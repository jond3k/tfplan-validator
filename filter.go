package tfplan_validator

import (
	"encoding/json"
	"io/ioutil"

	tfjson "github.com/hashicorp/terraform-json"
)

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

// FilterForPlan creates a filter that accepts everything in the specified plan
func NewFilterFromPlan(plan *tfjson.Plan) (*PlanFilter, error) {
	return &PlanFilter{}, nil
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
