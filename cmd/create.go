package cmd

import (
	"errors"
	"fmt"

	tfpv "github.com/fautom/tfplan-validator"
	"github.com/spf13/cobra"
)

func runCreateCmd(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("expected at least 2 arguments")
	}

	planPaths := args[0 : len(args)-1]
	outputPath := args[len(args)-1]

	if filter, err := tfpv.NewFilterFromPlanPaths(planPaths); err != nil {
		return err
	} else {
		filter.WriteJSON(outputPath)
		fmt.Printf("create %s %s", outputPath, planPaths)
	}
	return nil
}

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create PLAN_FILE... OUTPUT_FILE",
		Short: "Create a plan validator from one or more plans",
		RunE:  runCreateCmd,
	}
	return cmd
}
