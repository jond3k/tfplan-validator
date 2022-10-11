package cmd

import (
	"fmt"

	tfpv "github.com/fautom/tfplan-validator"
	"github.com/spf13/cobra"
)

func runPlanCmd(cmd *cobra.Command, args []string) error {
	var (
		globs      []string
		workspaces []string
		cacheDir   string
		command    string
		initArgs   string
		err        error
	)

	if globs, err = cmd.Flags().GetStringArray("glob"); err != nil {
		return err
	} else if cacheDir, err = cmd.Flags().GetString("cache-dir"); err != nil {
		return err
	} else if command, err = cmd.Flags().GetString("command"); err != nil {
		return err
	} else if initArgs, err = cmd.Flags().GetString("init-args"); err != nil {
		return err
	} else if workspaces, err = tfpv.FindWorkspaces(globs); err != nil {
		return err
	}

	fmt.Printf("Found %d workspaces %s\n", len(workspaces), workspaces)

	for _, workspace := range workspaces {
		if exec, err := tfpv.NewTerraformExec(command, initArgs, cacheDir, workspace); err != nil {
			return err
		} else if err = exec.Plan(); err != nil {
			return err
		}
	}

	return nil
}

func newPlanCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "plan",
		Short: "Runs plans on any workspaces it discovers in the current working directory and then creates a validator",
		RunE:  runPlanCmd,
	}
	cmd.Flags().String("cache-dir", ".tfpv-cache", "The workspace directory for collecting plans and rules")
	cmd.Flags().StringArrayP("glob", "g", tfpv.DefaultGlobs, "One or more globs to find terraform workspaces. Can use double-star wildcards and negation with !")
	cmd.Flags().StringP("command", "c", "", "The terraform command to use. By default it will use 'terragrunt' if there is a terragrunt.hcl file or 'terraform' otherwise")
	cmd.Flags().StringP("init-args", "i", "", "A string that contains additional args to pass to terraform init")
	return cmd
}
