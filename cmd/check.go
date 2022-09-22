package cmd

import (
	"errors"
	"fmt"
	"strings"

	tfpv "github.com/fautom/tfplan-validator"
	"github.com/spf13/cobra"
)

func printCheckReject(cmd *cobra.Command, results map[string]*tfpv.FilterResults) {
	out := cmd.ErrOrStderr()
	for path, result := range results {
		if result.HasErrors() {
			fmt.Fprintf(out, "The plan %s has been rejected because it has the following actions:\n\n", path)
			for addr, action := range result.Errors {
				if allowedActions := result.PlanFilter.AllowedActions[addr]; allowedActions == nil {
					fmt.Fprintf(out, "  - %s cannot be %s\n", addr, action.Pretty())
				} else {
					allowedPretty := make([]string, len(allowedActions))
					for i, allowedAction := range allowedActions {
						allowedPretty[i] = allowedAction.Pretty()
					}
					fmt.Fprintf(out, "  - %s cannot be %s only %s\n", addr, action.Pretty(), strings.Join(allowedPretty, " or "))
				}
			}
		}
		fmt.Fprint(out, "\n")
	}
}

func printCheckAccept(cmd *cobra.Command, results map[string]*tfpv.FilterResults) {
	out := cmd.OutOrStdout()
	for path, result := range results {
		if result.HasChanges() {
			fmt.Fprintf(out, "The plan %s passes checks and will perform the following actions:\n\n", path)
			for addr, action := range result.Changes {
				fmt.Fprintf(out, "  - %s will be %s\n", addr, action.Pretty())
			}
		}
	}
}

func runCheckCmd(cmd *cobra.Command, args []string) error {
	if len(args) < 1 {
		return errors.New("expected at least one plan")
	}

	planPaths := args
	rulesPath, err := cmd.Flags().GetString("rules")

	if err != nil {
		return err
	}

	plans, err := tfpv.ReadPlans(planPaths)

	if err != nil {
		return fmt.Errorf("failed to load plans: %w", err)
	}

	rules, err := tfpv.ReadPlanFilter(rulesPath)

	if err != nil {
		return fmt.Errorf("failed to read rules: %w", err)
	}

	results := map[string]*tfpv.FilterResults{}
	hasErrors := false

	for i, plan := range plans {
		result, err := tfpv.CheckPlan(rules, plan)
		if err != nil {
			return err
		}
		hasErrors = hasErrors || result.HasErrors()
		results[planPaths[i]] = result
	}

	if hasErrors {
		cmd.SilenceUsage = true
		printCheckReject(cmd, results)
		return errors.New("invalid plan")
	}

	printCheckAccept(cmd, results)

	return nil
}

func newCheckCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check PLAN_FILE... [--rules RULES_FILE]",
		Short: "Validate one or more plan using a rule file",
		RunE:  runCheckCmd,
	}
	cmd.Flags().String("rules", "./rules.json", "The rules file to use")
	return cmd
}
