package cmd

import (
	"testing"
)

func TestPlanCmd(t *testing.T) {
	cases := []cmdCase{
		{
			name:   "",
			args:   []string{"plan"},
			stdout: ``,
		},
	}

	for _, tc := range cases {
		tc.run(t)
	}

}
