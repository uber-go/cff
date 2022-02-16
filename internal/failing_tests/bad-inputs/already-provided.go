package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// AlreadyProvidedTaskParam is a function that provides a string type twice.
func AlreadyProvidedTaskParam() {
	cff.Flow(context.Background(),
		cff.Task(
			func() string {
				return "a"
			}),
		cff.Task(
			func() string {
				return "b"
			},
		),
	)
}

// AlreadyProvidedFlowParams is a function that provides multiple types multiple
// times to cff.Params.
func AlreadyProvidedFlowParams() {
	cff.Flow(
		context.Background(),
		cff.Params(1, 2, true, false),
		cff.Task(
			func(i int) int {
				return i
			},
		),
	)
}
