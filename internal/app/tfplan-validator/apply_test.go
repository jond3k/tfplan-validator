package cmd

import (
	"testing"
)

func TestApplyCmd(t *testing.T) {
	cases := []cmdCase{
		{
			name:   "",
			args:   []string{"apply"},
			stdout: ``,
		},
	}

	for _, tc := range cases {
		tc.run(t)
	}

}
