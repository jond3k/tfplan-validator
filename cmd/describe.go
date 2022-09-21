package cmd

import (
	"errors"
	"fmt"

	tfpv "github.com/fautom/tfplan-validator"
	"github.com/spf13/cobra"
)

func runDescribeCmd(cmd *cobra.Command, args []string) error {
	if len(args) != 1 {
		return errors.New("expected a path to a rule file")
	}
	filterPath := args[0]
	filter, err := tfpv.ReadPlanFilter(filterPath)

	if err != nil {
		return err
	}

	fmt.Fprintf(cmd.OutOrStdout(), "The rules file %s allows Terraform to perform the following actions:\n\n", filterPath)
	printPlanFilterLines(cmd, filterPath, filter)

	return nil
}

func newDescribeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "describe RULES_FILE",
		Short: "Pretty print the contents of a rules file",
		RunE:  runDescribeCmd,
	}
}
