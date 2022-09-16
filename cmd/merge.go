package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

func Merge(rulePaths []string) error {
	fmt.Printf("merge %s", rulePaths)
	return nil
}

func newMergeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "merge RULES_FILE RULES_FILE2...",
		Short: "Combine two or more rules files",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 2 {
				return errors.New("expected paths for at least 2 rule files")
			}
			return Merge(args)
		},
	}
}
