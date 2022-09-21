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
			name:   "failure",
			args:   []string{"check", "../fixtures/create/plan.json", "../fixtures/delete/filter.json"},
			stdout: ``,
			stderr: `The plan ../fixtures/create/plan.json has been rejected because it has the following actions:

  - local_file.foo cannot be created only deleted

Error: invalid plan`,
		},
	}

	for _, tc := range cases {
		tc.run(t)
	}

}
