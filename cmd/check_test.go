package cmd

import "testing"

func TestCheckCmd(t *testing.T) {
	cases := []cmdCase{
		{
			name: "success",
			args: []string{"check", "../fixtures/create/plan.json", "../fixtures/create/filter.json"},
			stdout: `The plan ../fixtures/create/plan.json passes checks and will perform the following actions:

  - local_file.foo will be created`,
		},
		{
			name:   "failure known resource",
			args:   []string{"check", "../fixtures/create/plan.json", "../fixtures/delete/filter.json"},
			stdout: ``,
			stderr: `The plan ../fixtures/create/plan.json has been rejected because it has the following actions:

  - local_file.foo cannot be created only deleted

Error: invalid plan`,
		},
		{
			name:   "unknown resource",
			args:   []string{"check", "../fixtures/create/plan.json", "../fixtures/update/filter.json"},
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
			args: []string{"check", "../fixtures/missing.json", "../fixtures/update/filter.json"},
			stdout: `Usage:
  tfplan-validator check PLAN_FILE... RULES_FILE [flags]

Flags:
  -h, --help   help for check`,
			stderr: `Error: failed to load plans: ../fixtures/missing.json: open ../fixtures/missing.json: no such file or directory`,
		},
		{
			name: "missing filter",
			args: []string{"check", "../fixtures/create/plan.json", "../fixtures/update/missing.json"},
			stdout: `Usage:
  tfplan-validator check PLAN_FILE... RULES_FILE [flags]

Flags:
  -h, --help   help for check`,
			stderr: `Error: failed to read rules: open ../fixtures/update/missing.json: no such file or directory`,
		},
		{
			name: "invalid plan",
			args: []string{"check", "../fixtures/itest/invalid-plan.json", "../fixtures/update/filter.json"},
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
