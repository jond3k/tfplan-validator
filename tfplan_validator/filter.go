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

// LoadPlanFilter from a path
func LoadPlanFilter(path string) (*PlanFilter, error) {
	if data, err := ioutil.ReadFile(path); err != nil {
		return nil, err
	} else {
		return ParsePlanFilter(data)
	}
}

// ParsePlanFilter from JSON
func ParsePlanFilter(data []byte) (*PlanFilter, error) {
	var f PlanFilter
	if err := json.Unmarshal(data, &f); err != nil {
		return nil, err
	}
	return &f, nil
}

// NewPlanFilterFromPlan creates a filter that accepts everything in the specified plan
func NewPlanFilterFromPlan(plan *tfjson.Plan) (*PlanFilter, error) {
	return &PlanFilter{}, nil
}

// MergePlanFilters combines filters from multiple plans into one
func MergePlanFilters([]*PlanFilter) (*PlanFilter, error) {
	return &PlanFilter{}, nil
}
