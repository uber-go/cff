// +build cff

package badinputs

import (
	"go.uber.org/cff"
	"context"
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
