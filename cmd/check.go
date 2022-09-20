package cmd

import (
	"errors"
	"fmt"

	tfpv "github.com/fautom/tfplan-validator"
	"github.com/spf13/cobra"
)

func runCheckCmd(cmd *cobra.Command, args []string) error {
	if len(args) < 2 {
		return errors.New("expected at least 2 arguments")
	}

	planPaths := args[0 : len(args)-1]
	rulesPath := args[len(args)-1]

	plans, err := tfpv.ReadPlans(planPaths)

	if err != nil {
		return fmt.Errorf("failed to load plans: %w", err)
	}

	rules, err := tfpv.ReadPlanFilter(rulesPath)

	if err != nil {
		return fmt.Errorf("failed to read rules: %w", err)
	}

	var results []*tfpv.FilterResults

	for _, plan := range plans {
		result, err := tfpv.CheckPlan(rules, plan)
		if err != nil {
			return err
		}
		results = append(results, result)
	}

	// TODO something with results!

	return nil
}

func newCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check PLAN_FILE... RULES_FILE",
		Short: "Validate one or more plan using a rule file",
		RunE:  runCheckCmd,
	}
}
