package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func runApplyCmd(cmd *cobra.Command, args []string) error {
	var (
		baseCacheDir string
		mf           *Manifest
		err          error
	)

	if baseCacheDir, err = cmd.Flags().GetString("cache-dir"); err != nil {
		return err
	} else if mf, err = LoadManifest(baseCacheDir); err != nil {
		return err
	}

	fmt.Printf("Loaded manifest with %d workspaces\n", len(mf.Workspaces))
	cmd.SilenceUsage = true

	if err = Apply(mf); err != nil {
		return err
	}

	return nil
}

func newApplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Runs apply on any workspaces it discovers in the current working directory with a validator",
		RunE:  runApplyCmd,
	}
	cmd.Flags().String("cache-dir", DefaultCacheDir, "The workspace directory for collecting plans and rules")
	return cmd
}
