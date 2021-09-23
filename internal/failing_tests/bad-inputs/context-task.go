package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// ContextTask is a flow whose task has a context argument that is not the first positional argument.
func ContextTask() {
	cff.Flow(context.Background(),
		cff.Task(
			func(string, context.Context) bool {
				return true
			},
		),
	)
}
