package cmd

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tfplan-validator",
		Short: "A simple way to validate Terraform plans. Designed to assist batch operations on large numbers of similar state files.",
	}
	rules := &cobra.Command{
		Use:   "rules",
		Short: "Subcommands for working with rules files",
	}
	rules.AddCommand(newCreateCmd())
	rules.AddCommand(newCheckCmd())
	rules.AddCommand(newDescribeCmd())
	rules.AddCommand(newMergeCmd())
	cmd.AddCommand(rules)
	return cmd
}
