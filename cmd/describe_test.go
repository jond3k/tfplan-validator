package cmd

import "testing"

func TestDescribeCmd(t *testing.T) {
	cases := []cmdCase{
		{
			name: "success",
			args: []string{"describe", "../fixtures/create/filter.json"},
			stdout: `The rules file ../fixtures/create/filter.json allows Terraform to perform the following actions:

  - local_file.foo can be created`,
		},
		{
			name: "missing file",
			args: []string{"describe", "../fixtures/create/missing.json"},
			stdout: `Usage:
  tfplan-validator describe RULES_FILE [flags]

Flags:
  -h, --help   help for describe`,
			stderr: `Error: open ../fixtures/create/missing.json: no such file or directory`,
		},
	}

	for _, tc := range cases {
		tc.run(t)
	}

}
