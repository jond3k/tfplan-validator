package cmd

import "testing"

func TestMergeCmd(t *testing.T) {
	cases := []cmdCase{
		{
			name: "success",
			args: []string{"rules", "merge", filterPath("create"), filterPath("delete-create"), resultPath("test-merge.json")},
			stdout: `Created rules file ` + resultPath("test-merge.json") + ` that allows Terraform to perform the following actions:

  -+ local_file.foo can be created or replaced (deleted then re-created)`,
			files: map[string]string{
				resultPath("test-merge.json"): loadTestData(otherPath("create-delete-create.json")),
			},
		},
		{
			name: "missing args",
			args: []string{"rules", "merge"},
			stdout: `Usage:
  tfplan-validator rules merge RULES_FILE... OUTPUT_FILE [flags]

Flags:
  -h, --help   help for merge`,
			stderr: `Error: expected paths for at least 2 rule files and an output path`,
		},
		{
			name: "missing rules",
			args: []string{"rules", "merge", filterPath("update"), filterPath("missing"), resultPath("test-merge.json")},
			stdout: `Usage:
  tfplan-validator rules merge RULES_FILE... OUTPUT_FILE [flags]

Flags:
  -h, --help   help for merge`,
			stderr: `Error: open ` + filterPath("missing") + `: no such file or directory`,
		},
		{
			name:   "reject contradition",
			args:   []string{"rules", "merge", filterPath("create"), filterPath("delete"), resultPath("test-merge.json")},
			stdout: ``,
			stderr: `Error: failed to merge filters: contradictory actions: local_file.foo has delete and create`,
		},
	}

	for _, tc := range cases {
		tc.run(t)
	}

}
