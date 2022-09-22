package cmd

import "testing"

func TestCreateCmd(t *testing.T) {
	cases := []cmdCase{
		{
			name: "sucess simple",
			args: []string{"create", planPath("create"), planPath("delete-create"), resultPath("test-create.json")},
			files: map[string]string{
				resultPath("test-create.json"): loadTestData(otherPath("create-delete-create.json")),
			},
			stdout: `Created rules file ` + resultPath("test-create.json") + ` that allows Terraform to perform the following actions:

  - local_file.foo can be created or replaced (deleted then re-created)`,
		},

		{
			name: "missing args",
			args: []string{"create"},
			stdout: `Usage:
  tfplan-validator create PLAN_FILE... OUTPUT_FILE [flags]

Flags:
  -h, --help   help for create`,
			stderr: `Error: expected at least 2 arguments`,
		},
	}

	for _, tc := range cases {
		tc.run(t)
	}

}
