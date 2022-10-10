package cmd

import "testing"

func TestCreateCmd(t *testing.T) {
	cases := []cmdCase{
		{
			name: "success simple",
			args: []string{"rules", "create", planPath("create"), planPath("delete-create"), "--rules", resultPath("test-create.json")},
			files: map[string]string{
				resultPath("test-create.json"): loadTestData(otherPath("create-delete-create.json")),
			},
			stdout: `Created rules file ` + resultPath("test-create.json") + ` that allows Terraform to perform the following actions:

  -+ local_file.foo can be created or replaced (deleted then re-created)`,
		},

		{
			name: "missing args",
			args: []string{"rules", "create"},
			stdout: `Usage:
  tfplan-validator rules create PLAN_FILE... [--rules RULES_FILE] [flags]

Flags:
  -h, --help           help for create
      --rules string   The rules file to write (default "./rules.json")`,
			stderr: `Error: expected at least one plan`,
		},
	}

	for _, tc := range cases {
		tc.run(t)
	}

}
