package cmd

import (
	"errors"
	"fmt"

	tfpv "github.com/fautom/tfplan-validator"
	"github.com/spf13/cobra"
)

func runMergeCmd(cmd *cobra.Command, args []string) error {
	if len(args) < 3 {
		return errors.New("expected paths for at least 2 rule files and an output path")
	}
	inputPaths := args[0 : len(args)-1]
	outputPath := args[len(args)-1]

	filters, err := tfpv.ReadPlanFilters(inputPaths)

	if err != nil {
		return err
	}

	merged, err := tfpv.MergePlanFilters(filters)

	if err != nil {
		cmd.SilenceUsage = true
		return fmt.Errorf("failed to merge filters: %w", err)
	}

	if err = merged.WriteJSON(outputPath); err != nil {
		return fmt.Errorf("failed to write json: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Created rules file %s that allows Terraform to perform the following actions:\n\n", outputPath)
	printPlanFilterLines(cmd, outputPath, merged)

	return nil
}

func newMergeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "merge RULES_FILE... OUTPUT_FILE",
		Short: "Combine two or more rules files and write them to a new path",
		RunE:  runMergeCmd,
	}
}
