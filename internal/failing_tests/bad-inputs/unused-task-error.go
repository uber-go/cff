//go:build cff && failing
// +build cff,failing

package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// NoInvokeWithError does not use cff.Invoke(true) on a task a single error result.
func NoInvokeWithError() {
	var out int8
	var s string
	cff.Flow(context.Background(),
		cff.Params(s),
		cff.Results(&out),
		cff.Task(func(s string) error {
			return nil
		}),
	)
}
