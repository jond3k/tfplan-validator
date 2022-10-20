package cmd

import (
	"github.com/spf13/cobra"
)

// DefaultCacheDir should be consistent between operations or the manifests wont be found
const DefaultCacheDir = ".tfpv-cache"

func newRulesCmd() *cobra.Command {
	rules := &cobra.Command{
		Use:   "rules",
		Short: "Subcommands for working with rules files",
	}
	rules.AddCommand(newCreateCmd())
	rules.AddCommand(newCheckCmd())
	rules.AddCommand(newDescribeCmd())
	rules.AddCommand(newMergeCmd())
	return rules
}

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tfplan-validator",
		Short: "A simple way to validate Terraform plans. Designed to assist batch operations on large numbers of similar state files.",
	}
	cmd.AddCommand(newRulesCmd())
	cmd.AddCommand(newPlanCmd())
	cmd.AddCommand(newApplyCmd())
	return cmd
}
