// +build cff

package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// AlreadyProvided is a function that provides a string type twice.
func AlreadyProvided() {
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
