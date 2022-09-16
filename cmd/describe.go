package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func Describe(rulesPath string) error {
	fmt.Printf("describe %s", rulesPath)
	return nil
}

func newDescribeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "describe RULES_FILE",
		Short: "Pretty print the contents of a rules file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return errors.New("expected a path to a rule file")
			}
			return Describe(args[0])
		},
	}
}
