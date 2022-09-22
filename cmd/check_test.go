package cmd

import (
	"path"
	"testing"
)

var missingPlanPath = path.Join("..", "fixtures", "missing.json")
var createPlanPath = path.Join("..", "fixtures", "create", "plan.json")
var createFilterPath = path.Join("..", "fixtures", "create", "filter.json")
var deleteFilterPath = path.Join("..", "fixtures", "delete", "filter.json")
var updateFilterPath = path.Join("..", "fixtures", "update", "filter.json")

func TestCheckCmd(t *testing.T) {
	cases := []cmdCase{
		{
			name: "success",
			args: []string{"check", createPlanPath, createFilterPath},
			stdout: `The plan ../fixtures/create/plan.json passes checks and will perform the following actions:

  - local_file.foo will be created`,
		},
		{
			name:   "failure known resource",
			args:   []string{"check", createPlanPath, deleteFilterPath},
			stdout: ``,
			stderr: `The plan ../fixtures/create/plan.json has been rejected because it has the following actions:

  - local_file.foo cannot be created only deleted

Error: invalid plan`,
		},
		{
			name:   "unknown resource",
			args:   []string{"check", createPlanPath, updateFilterPath},
			stdout: ``,
			stderr: `The plan ../fixtures/create/plan.json has been rejected because it has the following actions:

  - local_file.foo cannot be created

Error: invalid plan`,
		},
		{
			name: "missing args",
			args: []string{"check"},
			stdout: `Usage:
  tfplan-validator check PLAN_FILE... RULES_FILE [flags]

Flags:
  -h, --help   help for check`,
			stderr: `Error: expected at least 2 arguments`,
		},
		{
			name: "missing plan",
			args: []string{"check", missingPlanPath, updateFilterPath},
			stdout: `Usage:
  tfplan-validator check PLAN_FILE... RULES_FILE [flags]

Flags:
  -h, --help   help for check`,
			stderr: "Error: failed to load plans: " + missingPlanPath + ": open " + missingPlanPath + ": no such file or directory",
		},
		{
			name: "missing filter",
			args: []string{"check", createPlanPath, missingPlanPath},
			stdout: `Usage:
  tfplan-validator check PLAN_FILE... RULES_FILE [flags]

Flags:
  -h, --help   help for check`,
			stderr: `Error: failed to read rules: open ../fixtures/missing.json: no such file or directory`,
		},
		{
			name: "invalid plan",
			args: []string{"check", "../fixtures/itest/invalid-plan.json", updateFilterPath},
			stdout: `Usage:
  tfplan-validator check PLAN_FILE... RULES_FILE [flags]

Flags:
  -h, --help   help for check`,
			stderr: `Error: unrecognized action in plan: [invalid]`,
		},
	}

	for _, tc := range cases {
		tc.run(t)
	}

}
