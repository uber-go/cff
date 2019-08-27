package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// UnsupportedInvoke uses cff.Invoke(true) on a task with output.
func UnsupportedInvoke() {
	var out int8
	var s string
	cff.Flow(context.Background(),
		cff.Params(s),
		cff.Results(&out),
		cff.Task(func(s string) int8 {
			return int8(0)
		},
			cff.Invoke(true)),
	)
}

// NoInvokeNoResults does not use cff.Invoke(true) on a task with no results.
func NoInvokeNoResults() {
	var out int8
	var s string
	cff.Flow(context.Background(),
		cff.Params(s),
		cff.Results(&out),
		cff.Task(func(s string) {
			return
		}))
}
