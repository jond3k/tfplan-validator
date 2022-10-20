package cmd

import (
	"os"
	"testing"
)

func TestPlanCmd(t *testing.T) {
	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	cases := []cmdCase{
		{
			name: "",
			args: []string{"plan"},
			stdout: `Usage:
  tfplan-validator plan [searchPath]... [flags]

Flags:
      --cache-dir string   The workspace directory for collecting plans and rules (default ".tfpv-cache")
  -c, --command string     The terraform command to use. By default it will use 'terragrunt' if there is a terragrunt.hcl file or 'terraform' otherwise
  -g, --glob stringArray   One or more globs to find terraform workspaces. Can use double-star wildcards and negation with ! (default [**/*.tf,**/.terraform.lock.hcl,!**/modules/**/*.tf,!**/modules/**/.terraform.lock.hcl,!**/.terragrunt-cache/**/*.tf,!**/.terragrunt-cache/**/.terraform.lock.hcl])
  -h, --help               help for plan
  -i, --init-args string   A string that contains additional args to pass to terraform init`,
			stderr: `Error: unable to find workspaces in [` + wd + `] using globs [**/*.tf **/.terraform.lock.hcl !**/modules/**/*.tf !**/modules/**/.terraform.lock.hcl !**/.terragrunt-cache/**/*.tf !**/.terragrunt-cache/**/.terraform.lock.hcl]`,
		},
	}

	for _, tc := range cases {
		tc.run(t)
	}

}
