package cmd

import (
	"strings"

	tfpv "github.com/fautom/tfplan-validator"
)

func formatSymbolsForActions(actions []tfpv.Action) string {
	d, c, u := "", "", ""
	for _, action := range actions {
		switch action {
		case tfpv.ActionCreate:
			c = "+"
		case tfpv.ActionUpdate:
			u = "~"
		case tfpv.ActionDelete:
			d = "-"
		case tfpv.ActionDestroyBeforeCreate:
			d, c = "-", "+"
		case tfpv.ActionCreateBeforeDestroy:
			d, c = "-", "+"
		}
	}
	return strings.Join([]string{d, c, u}, "")
}
