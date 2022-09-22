package cmd

import (
	"github.com/spf13/cobra"
)

func New() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tfplan-validator",
		Short: "A simple way to validate Terraform plans. Designed to assist batch operations on large numbers of similar state files.",
	}
	cmd.AddCommand(newCreateCmd())
	cmd.AddCommand(newCheckCmd())
	cmd.AddCommand(newDescribeCmd())
	cmd.AddCommand(newMergeCmd())
	return cmd
}
