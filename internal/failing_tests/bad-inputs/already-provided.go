// +build cff

package badinputs

import (
	"go.uber.org/cff"
	"context"
)

// AlreadyProvided is a function that provides a string type twice
func AlreadyProvided() {
	cff.Flow(context.Background(),
		cff.Tasks(
			func() string {
				return "a"
			},
			func() string {
				return "b"
			},
		),
	)
}
