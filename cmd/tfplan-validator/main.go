package main

import (
	"os"

	cmd "github.com/fautom/tfplan-validator/internal/app/tfplan-validator"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
