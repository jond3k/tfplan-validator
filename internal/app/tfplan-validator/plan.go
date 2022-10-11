package cmd

import (
	"fmt"

	tfpv "github.com/fautom/tfplan-validator"
	"github.com/spf13/cobra"
)

func runPlanCmd(cmd *cobra.Command, args []string) error {
	var (
		globs         []string
		workspaceDirs []string
		baseCacheDir  string
		command       string
		initArgs      string
		mf            *Manifest
		err           error
	)

	if globs, err = cmd.Flags().GetStringArray("glob"); err != nil {
		return err
	} else if baseCacheDir, err = cmd.Flags().GetString("cache-dir"); err != nil {
		return err
	} else if command, err = cmd.Flags().GetString("command"); err != nil {
		return err
	} else if initArgs, err = cmd.Flags().GetString("init-args"); err != nil {
		return err
	} else if workspaceDirs, err = tfpv.FindWorkspaces(globs); err != nil {
		return err
	}

	if len(workspaceDirs) < 1 {
		return fmt.Errorf("unable to find workspaces using glob %s", globs)
	}

	fmt.Printf("Found %d workspaces %s\n", len(workspaceDirs), workspaceDirs)
	cmd.SilenceUsage = true

	if mf, err = NewManifest(command, initArgs, baseCacheDir, workspaceDirs); err != nil {
		return err
	} else if err = Plan(mf); err != nil {
		return err
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
