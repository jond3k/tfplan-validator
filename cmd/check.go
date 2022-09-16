package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func Check(rulesPath string, planPaths []string) error {
	fmt.Printf("check %s %s", rulesPath, planPaths)
	return nil
}

func newCheckCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "check RULES_FILE PLAN_FILE...",
		Short: "Validate one or more plan using a rule file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("expected paths for a rules file and one or more plans")
			}
			return Check(args[0], args[1:])
		},
	}
}
