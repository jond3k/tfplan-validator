package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func Create(rulesPath string, planPaths []string) error {
	fmt.Printf("create %s %s", rulesPath, planPaths)
	return nil
}

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create RULES_FILE PLAN_FILE...",
		Short: "Using one or more plan files create a rule file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("expected paths to the rules and plan files")
			}
			return Create(args[0], args[1:])
		},
	}
	return cmd
}
