package tfplan_validator

import (
	"fmt"
	"io/ioutil"

	tfjson "github.com/hashicorp/terraform-json"
)

func ReadPlans(paths []string) ([]*tfjson.Plan, error) {
	plans := make([]*tfjson.Plan, len(paths))
	for i, p := range paths {
		if plan, err := ReadPlan(p); err != nil {
			return nil, err
		} else {
			plans[i] = plan
		}
	}
	return plans, nil
}

func ReadPlan(path string) (*tfjson.Plan, error) {
	if data, err := ioutil.ReadFile(path); err != nil {
		return nil, err
	} else if plan, err := ParsePlan(data); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	} else {
		return plan, nil
	}
}

func ParsePlan(data []byte) (*tfjson.Plan, error) {
	plan := tfjson.Plan{}
	if err := plan.UnmarshalJSON(data); err != nil {
		return nil, err
	} else if err := plan.Validate(); err != nil {
		return nil, err
	}
	return &plan, nil
}
