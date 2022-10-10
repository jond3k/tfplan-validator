package cmd

import "testing"

func TestDescribeCmd(t *testing.T) {
	cases := []cmdCase{
		{
			name: "success",
			args: []string{"rules", "describe", filterPath("create")},
			stdout: `The rules file ` + filterPath("create") + ` allows Terraform to perform the following actions:

  + local_file.foo can be created`,
		},
		{
			name: "missing file",
			args: []string{"rules", "describe", filterPath("missing")},
			stdout: `Usage:
  tfplan-validator rules describe RULES_FILE [flags]

Flags:
  -h, --help   help for describe`,
			stderr: `Error: open ` + filterPath("missing") + `: no such file or directory`,
		},
	}

	for _, tc := range cases {
		tc.run(t)
	}

}
