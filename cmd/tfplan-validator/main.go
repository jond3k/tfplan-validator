package main

import (
	"os"

	"github.com/fautom/tfplan-validator/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		os.Exit(1)
	}
}
