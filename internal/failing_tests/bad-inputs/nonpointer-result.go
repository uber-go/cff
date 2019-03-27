// +build cff

package badinputs

import (
	"go.uber.org/cff"
	"context"
)

// ResultsNonPointer is a flow where we pass a variable that is not a pointer to cff.Results
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
