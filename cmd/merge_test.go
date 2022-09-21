package cmd

import "testing"

func TestMergeCmd(t *testing.T) {
	cases := []cmdCase{
		{
			name: "success",
			args: []string{"merge", "../fixtures/create/filter.json", "../fixtures/delete-create/filter.json", "../test-results/test-merge.json"},
			stdout: `Created rules file ../test-results/test-merge.json that allows Terraform to perform the following actions:

  - local_file.foo can be created or replaced (deleted then re-created)`,
			files: map[string]string{
				"../test-results/test-merge.json": loadTestData("../fixtures/itest/create-delete-create.json"),
			},
		},
		{
			name: "missing args",
			args: []string{"merge"},
			stdout: `Usage:
  tfplan-validator merge RULES_FILE... OUTPUT_FILE [flags]

Flags:
  -h, --help   help for merge`,
			stderr: `Error: expected paths for at least 2 rule files and an output path`,
		},
		{
			name: "missing rules",
			args: []string{"merge", "../fixtures/update/filter.json", "../fixtures/create/missing.json", "../test-results/test-merge.json"},
			stdout: `Usage:
  tfplan-validator merge RULES_FILE... OUTPUT_FILE [flags]

Flags:
  -h, --help   help for merge`,
			stderr: `Error: ../fixtures/create/missing.json: open ../fixtures/create/missing.json: no such file or directory`,
		},
		{
			name:   "reject contradition",
			args:   []string{"merge", "../fixtures/create/filter.json", "../fixtures/delete/filter.json", "../test-results/test-merge.json"},
			stdout: ``,
			stderr: `Error: failed to merge filters: contradictory actions: local_file.foo has delete and create`,
		},
	}

	for _, tc := range cases {
		tc.run(t)
	}

}
