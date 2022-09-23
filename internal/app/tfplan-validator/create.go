package cmd

import (
	"errors"
	"fmt"
	"strings"

	tfpv "github.com/fautom/tfplan-validator"
	"github.com/spf13/cobra"
)

func printPlanFilterLines(cmd *cobra.Command, rulesPath string, filter *tfpv.PlanFilter) {
	out := cmd.OutOrStdout()
	for addr, actions := range filter.AllowedActions {
		pretty := make([]string, len(actions))
		symbol := formatSymbolsForActions(actions)
		for i, action := range actions {
			pretty[i] = action.Pretty()
		}
		fmt.Fprintf(out, "  %s %s can be %s\n", symbol, addr, strings.Join(pretty, " or "))
	}
}

func runCreateCmd(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("expected at least one plan")
	}

	planPaths := args
	rulesPath, err := cmd.Flags().GetString("rules")

	if err != nil {
		return err
	}

	filter, err := tfpv.NewFilterFromPlanPaths(planPaths)

	if err != nil {
		return fmt.Errorf("failed to create filter from plans: %w", err)
	}

	if err = filter.WriteJSON(rulesPath); err != nil {
		return fmt.Errorf("failed to write json: %w", err)
	}

	fmt.Fprintf(cmd.OutOrStdout(), "Created rules file %s that allows Terraform to perform the following actions:\n\n", rulesPath)
	printPlanFilterLines(cmd, rulesPath, filter)

	return nil
}

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create PLAN_FILE... [--rules RULES_FILE]",
		Short: "Create a plan validator from one or more plans",
		RunE:  runCreateCmd,
	}
	cmd.Flags().String("rules", "./rules.json", "The rules file to write")
	return cmd
}
