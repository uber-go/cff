package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// ResultsNonPointer is a flow where we pass a variable that is not a pointer to cff.Results.
func ResultsNonPointer() {
	var v bool
	cff.Flow(context.Background(),
		cff.Results(v),
		cff.Task(
			func() bool {
				return true
			},
		),
	)
}
