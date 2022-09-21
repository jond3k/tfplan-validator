package cmd

import (
	"errors"
	"fmt"
	"strings"

	tfpv "github.com/fautom/tfplan-validator"
	"github.com/spf13/cobra"
)

func printPlanFilterLines(cmd *cobra.Command, filterPath string, filter *tfpv.PlanFilter) {
	out := cmd.OutOrStdout()
	for addr, actions := range filter.AllowedActions {
		pretty := make([]string, len(actions))
		for i, action := range actions {
			pretty[i] = action.Pretty()
		}
		fmt.Fprintf(out, "  - %s can be %s\n", addr, strings.Join(pretty, " or "))
	}
}

func runCreateCmd(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("expected at least 2 arguments")
	}

	planPaths := args[0 : len(args)-1]
	filterPath := args[len(args)-1]
	filter, err := tfpv.NewFilterFromPlanPaths(planPaths)

	if err != nil {
		return fmt.Errorf("failed to create filter from plans: %w", err)
	}

	if err = filter.WriteJSON(filterPath); err != nil {
		return fmt.Errorf("failed to write json: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Created rules file %s that allows Terraform to perform the following actions:\n\n", filterPath)
	printPlanFilterLines(cmd, filterPath, filter)

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
