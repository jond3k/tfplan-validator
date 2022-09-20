package tfplan_validator

import "fmt"

// MergePlanFilters combines filters from multiple plans into one
func MergePlanFilters(filters []*PlanFilter) (*PlanFilter, error) {
	allowed := map[Address][]Action{}

	for _, filter := range filters {
		for address, additions := range filter.AllowedActions {
			current := allowed[address]

			if current == nil {
				current = []Action{}
			}

			for _, l := range additions {
				duplicate := false
				for _, r := range current {
					if l == r {
						duplicate = true
					} else if !AreCompatible(l, r) {
						return nil, fmt.Errorf("contradictory actions: %s has %s and %s", address, l, r)
					}
				}
				if !duplicate {
					current = append(current, l)
				}
			}

			allowed[address] = current
		}
	}

	return &PlanFilter{
		FormatVersion:  CurrentFormatVersion,
		AllowedActions: allowed,
	}, nil
}
