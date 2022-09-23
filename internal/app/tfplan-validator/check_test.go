package cmd

import (
	"testing"
)

func TestCheckCmd(t *testing.T) {
	cases := []cmdCase{
		{
			name: "create",
			args: []string{"check", planPath("create"), "--rules", filterPath("create")},
			stdout: `The plan ` + planPath("create") + ` passes checks and will perform the following actions:

  + local_file.foo will be created`,
		},
		{
			name: "delete",
			args: []string{"check", planPath("delete"), "--rules", filterPath("delete")},
			stdout: `The plan ` + planPath("delete") + ` passes checks and will perform the following actions:

  - local_file.foo will be deleted`,
		},
		{
			name: "update",
			args: []string{"check", planPath("update"), "--rules", filterPath("update")},
			stdout: `The plan ` + planPath("update") + ` passes checks and will perform the following actions:

  ~ google_project_iam_policy.project will be updated`,
		},
		{
			name: "delete-create",
			args: []string{"check", planPath("delete-create"), "--rules", filterPath("delete-create")},
			stdout: `The plan ` + planPath("delete-create") + ` passes checks and will perform the following actions:

  -+ local_file.foo will be replaced (deleted then re-created)`,
		},
		{
			name: "create-delete",
			args: []string{"check", planPath("create-delete"), "--rules", filterPath("create-delete")},
			stdout: `The plan ` + planPath("create-delete") + ` passes checks and will perform the following actions:

  -+ local_file.foo will be replaced (re-created before deletion)`,
		},
		{
			name:   "failure known resource",
			args:   []string{"check", planPath("create"), "--rules", filterPath("delete")},
			stdout: ``,
			stderr: `The plan ` + planPath("create") + ` has been rejected because it has the following actions:

  ! local_file.foo cannot be created only deleted

Error: invalid plan`,
		},
		{
			name:   "unknown resource",
			args:   []string{"check", planPath("create"), "--rules", filterPath("update")},
			stdout: ``,
			stderr: `The plan ` + planPath("create") + ` has been rejected because it has the following actions:

  ! local_file.foo cannot be created

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
