package cmd

import (
	tfplan_validator "github.com/fautom/tfplan-validator"
	"github.com/spf13/cobra"
)

func runPlanCmd(cmd *cobra.Command, args []string) error {
	return nil
}

func newPlanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Runs plans on any workspaces it discovers in the current working directory and then creates a validator",
		RunE:  runPlanCmd,
	}
	cmd.Flags().String("cache-dir", ".tfpv-cache", "The workspace directory for collecting plans and rules")
	cmd.Flags().StringArrayP("workspace", "w", tfplan_validator.DefaultGlobs, "One or more globs to find terraform workspaces. Can use double-star wildcards and negation with !")
	cmd.Flags().StringP("command", "c", "", "The terraform command to use. By default it will use 'terragrunt' if there is a terragrunt.hcl file or 'terraform' otherwise")
	cmd.Flags().StringP("init-args", "i", "", "A string that contains additional args to pass to terraform init")
	return cmd
}
