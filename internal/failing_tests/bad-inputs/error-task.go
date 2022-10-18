//go:build cff && failing
// +build cff,failing

package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// TaskReturnsError is a flow whose error return value is not the last positional argument.
func TaskReturnsError() {
	cff.Flow(context.Background(),
		cff.Task(
			func() (error, bool) {
				return nil, true
			},
		),
	)
}
