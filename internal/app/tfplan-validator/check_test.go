package cmd

import (
	"testing"
)

func TestCheckCmd(t *testing.T) {
	cases := []cmdCase{
		{
			name: "success",
			args: []string{"check", planPath("create"), "--rules", filterPath("create")},
			stdout: `The plan ` + planPath("create") + ` passes checks and will perform the following actions:

  - local_file.foo will be created`,
		},
		{
			name:   "failure known resource",
			args:   []string{"check", planPath("create"), "--rules", filterPath("delete")},
			stdout: ``,
			stderr: `The plan ` + planPath("create") + ` has been rejected because it has the following actions:

  - local_file.foo cannot be created only deleted

Error: invalid plan`,
		},
		{
			name:   "unknown resource",
			args:   []string{"check", planPath("create"), "--rules", filterPath("update")},
			stdout: ``,
			stderr: `The plan ` + planPath("create") + ` has been rejected because it has the following actions:

  - local_file.foo cannot be created

Error: invalid plan`,
		},
		{
			name: "missing args",
			args: []string{"check"},
			stdout: `Usage:
  tfplan-validator check PLAN_FILE... --rules RULES_FILE [flags]

Flags:
  -h, --help           help for check
      --rules string   The rules file to use`,
			stderr: `Error: expected at least one plan`,
		},
		{
			name: "missing plan",
			args: []string{"check", planPath("missing"), "--rules", filterPath("update")},
			stdout: `Usage:
  tfplan-validator check PLAN_FILE... --rules RULES_FILE [flags]

Flags:
  -h, --help           help for check
      --rules string   The rules file to use`,
			stderr: "Error: failed to load plans: open " + planPath("missing") + ": no such file or directory",
		},
		{
			name: "missing filter",
			args: []string{"check", planPath("create"), "--rules", filterPath("missing")},
			stdout: `Usage:
  tfplan-validator check PLAN_FILE... --rules RULES_FILE [flags]

Flags:
  -h, --help           help for check
      --rules string   The rules file to use`,
			stderr: `Error: failed to read rules: open ` + filterPath("missing") + `: no such file or directory`,
		},
		{
			name: "invalid plan",
			args: []string{"check", otherPath("plan-invalid-actions.json"), "--rules", filterPath("update")},
			stdout: `Usage:
  tfplan-validator check PLAN_FILE... --rules RULES_FILE [flags]

Flags:
  -h, --help           help for check
      --rules string   The rules file to use`,
			stderr: `Error: unrecognized action in plan: [invalid]`,
		},
	}

	for _, tc := range cases {
		tc.run(t)
	}

}
